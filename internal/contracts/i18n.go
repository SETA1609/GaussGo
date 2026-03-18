package contracts

type I18nResolver interface {
	Resolve(modID string, locale string, key string) string
	ResolveWithOptions(modID string, locale string, key string, options map[string]any) string
	SupportedLocales() []string
}
