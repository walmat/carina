package profiles

type ProfileGroup struct {
	Name     string                 `json:"name"`
	Profiles map[string]ProfileData `json:"profiles"`
}

type ProfileData struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Billing      Address `json:"billing"`
	Shipping     Address `json:"shipping"`
	Payment      Payment `json:"payment"`
	MaxCheckouts int     `json:"maxCheckouts"`
}

type AbbrevPair struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Address struct {
	Name     string      `json:"name"`
	Line1    string      `json:"line1"`
	Line2    string      `json:"line2"`
	Line3    string      `json:"line3"`
	PostCode string      `json:"postCode"`
	City     string      `json:"city"`
	Country  AbbrevPair  `json:"country"`
	State    *AbbrevPair `json:"state"`
}

type Payment struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Type     string `json:"type"`
	Number   string `json:"number"`
	ExpMonth string `json:"expMonth"`
	ExpYear  string `json:"expYear"`
	Cvv      string `json:"cvv"`
}
