package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func containsTag(tags []string, tag string) bool {
	for _, s := range tags {
		if s == tag {
			return true
		}
	}
	return false
}

func getThemesByTag(themeList []map[string]interface{}, tag string) []map[string]interface{} {
	var themes []map[string]interface{}
	for _, item := range themeList {
		tTags, ok := item["Tags:"].([]string)
		if ok {
			if containsTag(tTags, tag) {
				themes = append(themes, item)
			}
		}
	}
	return themes
}

func sortByKey(themeList []map[string]interface{}, key string) {
	sort.Slice(themeList, func(i, j int) bool {
		key1, _ := themeList[i][key].(string)
		key2, _ := themeList[j][key].(string)
		num1, _ := strconv.Atoi(key1)
		num2, _ := strconv.Atoi(key2)
		return num1 > num2
	})
}

func getOrDefault(value interface{}, defaultValue string) string {
	if value != nil {
		return fmt.Sprintf("%v", value)
	}
	return defaultValue
}

func matchThemeName(url any) string {
	sUrl, _ := url.(string)
	re := regexp.MustCompile(`/themes/([^/]+)/$`)
	return re.FindStringSubmatch(sUrl)[1]
}

func generateContent(tag string, themes []map[string]interface{}) string {
	var content string
	defaultValue := "unknown"
	content += "<details>\n"
	content += "<summary>" + tag + "</summary>\n <br>\n"
	content += " \n"
	content += "| Theme | Author | Stars | Updated | Minimum Hugo Version | License |\n"
	content += "| ------- | ------- | ------- | ------- | ------- | ------- |\n"
	for _, theme := range themes {
		line := fmt.Sprintf("| [%s](%s) | %s | %s | %s | %s | %s |\n",
			matchThemeName(theme["url"]),
			getOrDefault(theme["url"], defaultValue),
			getOrDefault(theme["Author:"], defaultValue),
			getOrDefault(theme["GitHub Stars:"], defaultValue),
			getOrDefault(theme["Updated:"], defaultValue),
			getOrDefault(theme["Minimum Hugo Version:"], defaultValue),
			getOrDefault(theme["License:"], defaultValue),
		)
		content += line
	}
	content += "</details>\n"
	return content
}

func generateREADME(content string) {
	err := os.WriteFile("test.md", []byte(content), 0644)
	if err != nil {
		log.Println("generateREADME error!")
		return
	}
	log.Println("generateREADME success!")
}
