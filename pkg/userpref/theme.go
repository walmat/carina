package userpref

var (
	themeState uint8
)

func SetThemeState(i uint8) {
	themeState = i
}

func GetThemeState() uint8 {
	return themeState
}
