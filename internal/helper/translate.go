package helper

type Locale string

func (l Locale) String() string {
	return string(l)
}

type Translates map[Locale]string

func (t Translates) Get(l Locale) string {
	return t[l]
}

func (t Translates) Set(l Locale, s string) {
	t[l] = s
}
