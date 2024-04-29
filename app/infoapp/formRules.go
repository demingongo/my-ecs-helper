package infoapp

func selectRuleJSON(info string, title string, dir string) string {
	value := selectJSONFile(title, dir, info)
	return value
}
