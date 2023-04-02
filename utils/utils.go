package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"

	"strings"
)

var pods map[string]interface{}

func ProcessHtml(BodyHtml string) string {
	xpath := "//h1[@itemprop='headline']/text()"
	fmt.Println(xpath)
	doc, err := htmlquery.Parse(strings.NewReader(BodyHtml))

	nodes, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		fmt.Println("Got Error: ", err)
	} else {
		fmt.Println("nodes: ", htmlquery.InnerText(nodes[0]))
	}

	for i, j := range nodes {
		dt := htmlquery.FindOne(j, "//text()")
		fmt.Println(i, htmlquery.InnerText(dt))
	}
	return xpath
}

func GetMyFunc() {
	/*
		These functions are executed in the following order:
		OnRequest(): Called before performing an HTTP request with Visit().
		OnError(): Called if an error occurred during the HTTP request.
		OnResponse(): Called after receiving a response from the server.
		OnHTML(): Called right after OnResponse() if the received content is HTML.
		OnScraped(): Called after all OnHTML() callback executions.
	*/
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
		// fmt.Println("body:", string(r.Body))
		// return string(r.Body)
		ProcessHtml(string(r.Body))
	})

	//c.OnHTML("body", func(e *colly.HTMLElement) {
	//	fmt.Println(e.Text)
	//})

	c.Visit("https://www.europapress.es/castilla-lamancha/noticia-lm-priorizara-empresas-incluyan-menus-productos-ecologicos-centros-sociosanitarios-20230214131303.html")
	//fmt.Println("I am from utils", c.Name)
	//fmt.Println("I am from utils", c.UrlPrifix)
}

func getAbsoluteURL(rawUrl string, rssUrl string, config string) string {
	storyUrl := rawUrl

	if config != "" {

		json.Unmarshal([]byte(config), &pods)
		if !((strings.HasPrefix(rawUrl, "http://")) || (strings.HasPrefix(rawUrl, "https://"))) {

			if pods["prefix"] != nil {
				storyUrl = pods["prefix"].(string) + rawUrl
			} else {
				u, _ := url.Parse(rssUrl)
				domailUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				if strings.HasPrefix(rawUrl, "?") {
					storyUrl = domailUrl + u.Path
				} else if strings.HasPrefix(rawUrl, "/") {
					storyUrl = domailUrl + rawUrl
				} else if strings.HasPrefix(rawUrl, "../") {
					rawUrl = strings.Replace(rawUrl, "..", "", 1)
					storyUrl = domailUrl + rawUrl
				} else if strings.HasPrefix(rawUrl, "//www.") || strings.HasPrefix(rawUrl, "//") {
					storyUrl = u.Scheme + ":" + rawUrl

				}
			}

		}

		breakPoint := pods["break_point"]
		removablePatterns := pods["removable_patterns"]
		suffix := pods["suffix"]
		if breakPoint != nil || removablePatterns != nil {
			fmt.Println("removable_patterns: ", breakPoint, removablePatterns)
			storyUrl = cleanStoryUrl(storyUrl, breakPoint.(string), removablePatterns.(string))
		}
		if suffix != nil {
			storyUrl = storyUrl + suffix.(string)
		}

	}
	return storyUrl
}

func cleanStoryUrl(url string, breakPoint string, removablePatterns string) string {
	cleanUrl := url
	fmt.Println("removable_patterns: ", pods["removable_patterns"], breakPoint, removablePatterns)
	if breakPoint != "" {
		cleanUrl = strings.Split(url, breakPoint)[0]
	} else if removablePatterns != "" {
		for _, regex := range removablePatterns {
			m := regexp.MustCompile(string(regex))
			url = m.ReplaceAllString(url, "")
		}
	}
	return cleanUrl
}
