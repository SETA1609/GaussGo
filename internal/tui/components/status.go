package components

func RenderStatusPill(theme Theme, text string) string {
	if text == "" {
		return ""
	}
	return theme.StatusPill.Render(text)
}
