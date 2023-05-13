package main

func main() {
	var tags = []string{"Blog", "Responsive", "Minimal", "Personal", "Light", "Dark", "Multilingual", "Portfolio", "Bootstrap", "Landing", "Dark Mode", "Docs", "Company", "Gallery", "Contact", "Archive"}
	hugoCrawler := GetCrawler()
	hugoCrawler.InitHugoThemeCrawler()
	err := hugoCrawler.startCrawlHugoThemes()
	if err == nil {
		list := hugoCrawler.getThemes()
		var content string
		for _, tag := range tags {
			temp := getThemesByTag(list, tag)
			sortByKey(temp, "GitHub Stars:")
			content += generateContent(tag, temp)
		}
		generateREADME(content)
	}
}
