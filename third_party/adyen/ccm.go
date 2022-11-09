package adyen

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"math"
)

// from github.com/pion/dtls/tree/master/pkg/crypto/ccm

type ccm struct {
	b cipher.Block
	M uint8
	L uint8
}

const ccmBlockSize = 16

type CCM interface {
	cipher.AEAD
	MaxLength() int
}

var (
	errInvalidBlockSize = errors.New("ccm: NewCCM requires 128-bit block cipher")
	errInvalidTagSize   = errors.New("ccm: tagsize must be 4, 6, 8, 10, 12, 14, or 16")
	errInvalidNonceSize = errors.New("ccm: invalid nonce size")
)

func NewCCM(b cipher.Block, tagSize, nonceSize int) (CCM, error) {
	if b.BlockSize() != ccmBlockSize {
		return nil, errInvalidBlockSize
	}
	if tagSize < 4 || tagSize > 16 || tagSize&1 != 0 {
		return nil, errInvalidTagSize
	}
	lensize := 15 - nonceSize
	if lensize < 2 || lensize > 8 {
		return nil, errInvalidNonceSize
	}
	c := &ccm{b: b, M: uint8(tagSize), L: uint8(lensize)}
	return c, nil
}

func (c *ccm) NonceSize() int { return 15 - int(c.L) }
func (c *ccm) Overhead() int  { return int(c.M) }
func (c *ccm) MaxLength() int { return maxLen(c.L, c.Overhead()) }

func maxLen(l uint8, tagSize int) int {
	max := (uint64(1) << (8 * l)) - 1
	if m64 := uint64(math.MaxInt64) - uint64(tagSize); l > 8 || max > m64 {
		max = m64
	}
	if max != uint64(int(max)) {
		return math.MaxInt32 - tagSize
	}
	return int(max)
}

func MaxNonceLength(dataLen int) int {
	const tagSize = 16
	for L := 2; L <= 8; L++ {
		if maxLen(uint8(L), tagSize) >= dataLen {
			return 15 - L
		}
	}
	return 0
}

func (c *ccm) cbcRound(mac, data []byte) {
	for i := 0; i < ccmBlockSize; i++ {
		mac[i] ^= data[i]
	}
	c.b.Encrypt(mac, mac)
}

func (c *ccm) cbcData(mac, data []byte) {
	for len(data) >= ccmBlockSize {
		c.cbcRound(mac, data[:ccmBlockSize])
		data = data[ccmBlockSize:]
	}
	if len(data) > 0 {
		var block [ccmBlockSize]byte
		copy(block[:], data)
		c.cbcRound(mac, block[:])
	}
}

var errPlaintextTooLong = errors.New("ccm: plaintext too large")

func (c *ccm) tag(nonce, plaintext, data []byte) ([]byte, error) {
	var mac [ccmBlockSize]byte

	if len(data) > 0 {
		mac[0] |= 1 << 6
	}
	mac[0] |= (c.M - 2) << 2
	mac[0] |= c.L - 1
	if len(nonce) != c.NonceSize() {
		return nil, errInvalidNonceSize
	}
	if len(plaintext) > c.MaxLength() {
		return nil, errPlaintextTooLong
	}
	binary.BigEndian.PutUint64(mac[ccmBlockSize-8:], uint64(len(plaintext)))
	copy(mac[1:ccmBlockSize-c.L], nonce)
	c.b.Encrypt(mac[:], mac[:])

	var block [ccmBlockSize]byte
	if n := uint64(len(data)); n > 0 {
		// First adata block includes adata length
		i := 2
		if n <= 0xfeff {
			binary.BigEndian.PutUint16(block[:i], uint16(n))
		} else {
			block[0] = 0xfe
			block[1] = 0xff
			if n < uint64(1<<32) {
				i = 2 + 4
				binary.BigEndian.PutUint32(block[2:i], uint32(n))
			} else {
				i = 2 + 8
				binary.BigEndian.PutUint64(block[2:i], n)
			}
		}
		i = copy(block[i:], data)
		c.cbcRound(mac[:], block[:])
		c.cbcData(mac[:], data[i:])
	}

	if len(plaintext) > 0 {
		c.cbcData(mac[:], plaintext)
	}

	return mac[:c.M], nil
}

func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func (c *ccm) Seal(dst, nonce, plaintext, data []byte) []byte {
	tag, err := c.tag(nonce, plaintext, data)
	if err != nil {
		panic(err)
	}

	var iv, s0 [ccmBlockSize]byte
	iv[0] = c.L - 1
	copy(iv[1:ccmBlockSize-c.L], nonce)
	c.b.Encrypt(s0[:], iv[:])
	for i := 0; i < int(c.M); i++ {
		tag[i] ^= s0[i]
	}
	iv[len(iv)-1] |= 1
	stream := cipher.NewCTR(c.b, iv[:])
	ret, out := sliceForAppend(dst, len(plaintext)+int(c.M))
	stream.XORKeyStream(out, plaintext)
	copy(out[len(plaintext):], tag)
	return ret
}

var (
	errOpen               = errors.New("ccm: message authentication failed")
	errCiphertextTooShort = errors.New("ccm: ciphertext too short")
	errCiphertextTooLong  = errors.New("ccm: ciphertext too long")
)

func (c *ccm) Open(dst, nonce, ciphertext, data []byte) ([]byte, error) {
	if len(ciphertext) < int(c.M) {
		return nil, errCiphertextTooShort
	}
	if len(ciphertext) > c.MaxLength()+c.Overhead() {
		return nil, errCiphertextTooLong
	}

	tag := make([]byte, int(c.M))
	copy(tag, ciphertext[len(ciphertext)-int(c.M):])
	ciphertextWithoutTag := ciphertext[:len(ciphertext)-int(c.M)]

	var iv, s0 [ccmBlockSize]byte
	iv[0] = c.L - 1
	copy(iv[1:ccmBlockSize-c.L], nonce)
	c.b.Encrypt(s0[:], iv[:])
	for i := 0; i < int(c.M); i++ {
		tag[i] ^= s0[i]
	}
	iv[len(iv)-1] |= 1
	stream := cipher.NewCTR(c.b, iv[:])

	plaintext := make([]byte, len(ciphertextWithoutTag))
	stream.XORKeyStream(plaintext, ciphertextWithoutTag)
	expectedTag, err := c.tag(nonce, plaintext, data)
	if err != nil {
		return nil, err
	}

	if subtle.ConstantTimeCompare(tag, expectedTag) != 1 {
		return nil, errOpen
	}
	return append(dst, plaintext...), nil
}
