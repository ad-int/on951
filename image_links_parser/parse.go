package image_links_parser

func findAllLinks(text string) []string {
	return []string{}
}

func storeImageLinks(links []string) map[string]string {
	return map[string]string{}
}

func updateImageLinks(text string, links map[string]string) string {
	return text
}

func Process(text string) (string, bool) {
	links := findAllLinks(text)
	storedLinks := storeImageLinks(links)
	return updateImageLinks(text, storedLinks), true

}
