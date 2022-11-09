package adyen2

func NewFingerprintString() string {
	fp := NewDF()
	fp.random()
	return fp.String()
}
