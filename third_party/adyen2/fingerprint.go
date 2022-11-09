package adyen2

import "regexp"

type Fingerprint struct {
	canvas        string
	cpuClass      string
	doNotTrack    string
	fonts         string
	jsFonts       string
	mimeTypes     string
	nrOfFonts     string
	nrOfMimeTypes string
	nrOfPlugins   string
	platform      string
	plugins       string
	superCookies  string
	timeZone      string
	userAgent     string
	video         string
	webglFp       string
	result        string
}

var rpPlus, _ = regexp.Compile(`\+`)
var rpSp, _ = regexp.Compile(`\/`)

func NewDF() *Fingerprint {
	fp := &Fingerprint{
		canvas:        "mvcujCW8QL",
		cpuClass:      "00000",
		doNotTrack:    "00000",
		fonts:         "0000000000",
		jsFonts:       "EC4FlSABmQ",
		mimeTypes:     "qBgAVdzoFe",
		nrOfFonts:     "000",
		nrOfMimeTypes: "004",
		nrOfPlugins:   "001",
		platform:      "qZkTE",
		plugins:       "cJRxL5CUm7",
		superCookies:  "cVB94iKzBG",
		timeZone:      "5B8+c9qXKS",
		userAgent:     "JnveW63V8k",
		video:         "0050271576",
		webglFp:       "Q6Fc6BeHpl",
	}
	fp.result = encodeAdyenFp(fp)
	return fp
}

//随机修改一些数据
func (fp *Fingerprint) random() {

}

//dfGetEntropy
func (fp *Fingerprint) dfGetEntropy() string {
	return "40"
}

//toString
func (fp *Fingerprint) String() string {
	return fp.result
}

func padStr(c string, a int) string {
	if len(c) >= a {
		return c[:a]
	}
	b := ""
	for ; len(b) < a-len(c); b += "0" {
	}
	return b
}

func encodeAdyenFp(a *Fingerprint) string {
	str := a.plugins + a.nrOfPlugins + a.fonts + a.nrOfFonts + a.timeZone + a.video + a.superCookies + a.userAgent + a.mimeTypes + a.nrOfMimeTypes + a.canvas + a.cpuClass + a.platform + a.doNotTrack + a.webglFp + a.jsFonts
	str = rpPlus.ReplaceAllString(str, "G")
	str = rpSp.ReplaceAllString(str, "D")
	return str + ":" + a.dfGetEntropy()
}
