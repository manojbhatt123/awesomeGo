package crawl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

type Column struct {
	web_url            interface{}
	last_run           interface{}
	rss_url            string
	stories_selector   interface{}
	url_selector       string
	title_selector     interface{}
	pub_date_selector  interface{}
	summary_selector   interface{}
	authors_selector   interface{}
	integration_config string
	req_data_j         interface{}
	req_data_t         interface{}
	req_headers        string
	id                 int64
	req_method         int
	proxy_type         interface{}
	archive_day_range  int
	category           int
	status             bool
}

type MyGeneric interface {
	Column | int | string
}

var (
	web_url            interface{}
	last_run           interface{}
	rss_url            string
	stories_selector   interface{}
	url_selector       string
	title_selector     interface{}
	pub_date_selector  interface{}
	summary_selector   interface{}
	authors_selector   interface{}
	integration_config interface{}
	req_data_j         interface{}
	req_data_t         interface{}
	req_headers        interface{}
	id                 int64
	req_method         int
	proxy_type         interface{}
	archive_day_range  int
	category           int
	status             bool
)

var pods map[string]interface{}

func CrawlUrl(col Column, finished chan bool) {
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		headers := col.req_headers
		json.Unmarshal([]byte(headers), &pods)
		for k, v := range pods {
			r.Headers.Set(k, v.(string))
		}
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
		data := ProcessHtml2(string(r.Body), &col)

		addOrUpdateRecord(data, r.StatusCode, &col)

		finished <- true
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
		finished <- true
	})

	c.Visit(col.rss_url)
}

func ProcessData(rows *sql.Rows) {
	var col Column
	finished := make(chan bool, 10)
	for rows.Next() {
		rows.Scan(&col.id, &col.rss_url, &col.stories_selector, &col.url_selector, &col.title_selector, &col.pub_date_selector, &col.summary_selector, &col.authors_selector, &col.archive_day_range, &col.req_headers)
		fmt.Println("Process Data", col.pub_date_selector)
		go CrawlUrl(col, finished)

	}

	for i := 0; i < 10; i++ {
		fmt.Println(<-finished)
	}
}

func ProcessHtml2(BodyHtml string, col *Column) []map[string]string {
	fmt.Println("$$$$$$$$ProcessHtml2$$$$$$$$$")
	doc, _ := htmlquery.Parse(strings.NewReader(BodyHtml))

	selectorXpathMap := make(map[string][]string, 10)

	if col.title_selector != nil {
		titles := []string{}
		title_selector := col.title_selector.(string)
		titleNodes, _ := htmlquery.QueryAll(doc, title_selector)
		for _, j := range titleNodes {
			dt := htmlquery.FindOne(j, "//text()")
			titles = append(titles, htmlquery.InnerText(dt))

		}
		selectorXpathMap["titles"] = titles
	}

	if col.pub_date_selector != nil {
		publishDates := []string{}
		pub_date_selector := col.pub_date_selector.(string)
		titleNodes, _ := htmlquery.QueryAll(doc, pub_date_selector)
		for _, j := range titleNodes {
			dt := htmlquery.FindOne(j, "//text()")
			publishDates = append(publishDates, htmlquery.InnerText(dt))
		}
		selectorXpathMap["pubDate"] = publishDates
	}

	if col.summary_selector != nil {
		summarys := []string{}
		summary_selector := col.summary_selector.(string)
		titleNodes, _ := htmlquery.QueryAll(doc, summary_selector)
		for _, j := range titleNodes {
			dt := htmlquery.FindOne(j, "//text()")
			summarys = append(summarys, htmlquery.InnerText(dt))
		}
		selectorXpathMap["summarys"] = summarys
	}

	if col.authors_selector != nil {
		authors := []string{}
		authors_selector := col.authors_selector.(string)
		titleNodes, _ := htmlquery.QueryAll(doc, authors_selector)
		for _, j := range titleNodes {
			dt := htmlquery.FindOne(j, "//text()")
			authors = append(authors, htmlquery.InnerText(dt))
		}
		selectorXpathMap["authors"] = authors
	}

	urls := []string{}
	urlSelector := col.url_selector
	urlNodes, _ := htmlquery.QueryAll(doc, urlSelector)
	for _, j := range urlNodes {
		dt := htmlquery.FindOne(j, "//text()")
		urls = append(urls, htmlquery.InnerText(dt))
	}
	selectorXpathMap["urls"] = urls

	var results []map[string]string
	urlSliceLength := len(selectorXpathMap["urls"])
	titleSliceLength := len(selectorXpathMap["titles"])
	pubDateSliceLength := len(selectorXpathMap["pubDate"])
	summarySliceLength := len(selectorXpathMap["summarys"])
	authorSliceLength := len(selectorXpathMap["authors"])
	for index, url := range selectorXpathMap["urls"] {
		items := map[string]string{}
		// url := getAbsoluteURL(url, col.rss_url, col.integration_config)
		items["url"] = url
		items["company_rss"] = col.rss_url

		if urlSliceLength == titleSliceLength {
			items["title"] = selectorXpathMap["titles"][index]
		} else if index >= titleSliceLength {
			items["title"] = ""
		} else if index < titleSliceLength {
			items["title"] = selectorXpathMap["titles"][index]
		}

		if urlSliceLength == pubDateSliceLength {
			items["pubDate"] = selectorXpathMap["pubDate"][index]
		} else if index >= pubDateSliceLength {
			items["pubDate"] = ""
		} else if index < pubDateSliceLength {
			items["pubDate"] = selectorXpathMap["pubDate"][index]
		}

		if urlSliceLength == summarySliceLength {
			items["summary"] = selectorXpathMap["summarys"][index]
		} else if index >= summarySliceLength {
			items["summary"] = ""
		} else if index < summarySliceLength {
			items["summary"] = selectorXpathMap["summarys"][index]
		}

		if urlSliceLength == authorSliceLength {
			items["author"] = selectorXpathMap["authors"][index]
		} else if index >= authorSliceLength {
			items["author"] = ""
		} else if index < authorSliceLength {
			items["author"] = selectorXpathMap["authors"][index]
		}

		results = append(results, items)
	}
	return results
}

func ProcessHtml(BodyHtml string, col *Column, finished chan bool) {

	doc, err := htmlquery.Parse(strings.NewReader(BodyHtml))
	urlNodes, err := htmlquery.QueryAll(doc, col.url_selector)
	if err != nil {
		fmt.Println("Got Error: ", err)
	}
	fmt.Println("urlNodes", col.id)
	fmt.Println("urlNodes", urlNodes)
	for i, j := range urlNodes {
		dt := htmlquery.FindOne(j, "//text()")
		fmt.Println(i, htmlquery.InnerText(dt))
	}
	fmt.Println("##########################hvbxjzhvbjh##################################")
	finished <- true

}

func ExecuteQuery(db *sql.DB, rawQuery string) *sql.Rows {
	rows, _ := db.Query(rawQuery)
	fmt.Println("Rows: ", rows)
	db.Close()
	return rows
}

func ExecuteWithoutGoroutineQuery(db *sql.DB, rawQuery string) {
	rows, _ := db.Query(rawQuery)
	db.Close()
	for rows.Next() {
		rows.Scan(&authors_selector, &summary_selector, &req_data_j, &req_data_t)
		fmt.Println(authors_selector, summary_selector, req_data_j, req_data_t)
	}
}
