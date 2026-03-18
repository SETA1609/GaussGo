package contracts

type I18nResolver interface {
	Resolve(modID string, locale string, key string) string
	SupportedLocales() []string
}
