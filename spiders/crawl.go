package crawl

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	// util "awesomeGo/utils"

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

// var urlCh chan int
// var htmlBodyData chan string

func CrawlUrl(url string, htmlBodyData chan string) {
	c := colly.NewCollector(
	//colly.Async(true),
	)
	c.SetRequestTimeout(10 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
		// fmt.Printf("Type Status: %T", r.StatusCode)
		// urlCh <- r.StatusCode
		// time.Sleep(1 * time.Second)
		fmt.Println("###########status##########", r.StatusCode)
		htmlBodyData <- string(r.Body)

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
		htmlBodyData <- string(r.Body)
	})

	c.Visit(url)
}

func ProcessData(rows *sql.Rows) {
	var col Column
	htmlBodyData := make(chan string, 3)
	finished := make(chan bool, 3)
	var wg sync.WaitGroup
	wg.Add(3)
	for rows.Next() {
		rows.Scan(&col.id, &col.rss_url, &col.stories_selector, &col.url_selector, &col.title_selector, &col.pub_date_selector, &col.summary_selector, &col.authors_selector, &col.proxy_type)

		go CrawlUrl(col.rss_url, htmlBodyData)
		RawHtml := <-htmlBodyData
		// fmt.Printf("Type %T", col)
		go ProcessHtml(RawHtml, &col, finished)
		// time.Sleep(1 * time.Second)
		fmt.Println("id: ", col.id)
		// fmt.Println(col.rss_url)
		// fmt.Println(col.stories_selector)
		// fmt.Println(col.url_selector)
		// fmt.Println(col.title_selector)

	}
	// for htmlDt := range htmlBodyData {
	// 	fmt.Println("for id: ", col.id)
	// 	ProcessHtml(htmlDt, col, finished, &wg)
	// 	fmt.Println("Done", <-finished)
	// }
	// wg.Wait()

	for i := 0; i < 3; i++ {
		fmt.Println(<-htmlBodyData)
	}
	close(htmlBodyData)
	close(finished)
	// for range htmlBodyData {
	// 	fmt.Println("for")
	// 	ProcessHtml(<-htmlBodyData)
	// }
}

// func (col2 Column) ProcessHtml2(htmlRawData string){
// 	fmt.Println(col2.id)
// }

func ProcessHtml2(BodyHtml string, col *Column) string {
	xpath := "//h1[@class='h1']/text()"
	fmt.Println("rssUrl: ", col.id)
	doc, err := htmlquery.Parse(strings.NewReader(BodyHtml))
	fmt.Println("doc", err)
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

func ProcessHtml(BodyHtml string, col *Column, finished chan bool) {
	//xpath := "//h1[@class='h1']/text()"
	// defer wg.Done()
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
	// return xpath
}

func ExecuteQuery(db *sql.DB, rawQuery string) *sql.Rows {
	// done := make(chan bool)
	rows, _ := db.Query(rawQuery)
	db.Close()
	return rows
	// col, _ := rows.ColumnTypes()

	// for rows.Next() {
	// c := Column{}
	// go c.ProcessData(rows)
	// type rssData interface {
	// 	last_run string;
	// }
	// var last_run, id, rss_url, status, category, archive_day_range, proxy_type, req_method, req_headers, req_data_j, req_data_t, integration_config, authors_selector, summary_selector, pub_date_selector, title_selector, url_selector, stories_selector, web_url interface{}
	// rows.Scan(Column.last_run, &id, &rss_url, &status, &category, &archive_day_range, &proxy_type, &req_method, &req_headers, &req_data_j, &req_data_t, &req_data_j, &integration_config, &authors_selector,
	// 	&summary_selector, &pub_date_selector, &title_selector, &url_selector, &stories_selector, &web_url)
	// go func(last_run string, id int64, rss_url string, status bool, category int, archive_day_range int, proxy_type interface{}, req_method int, req_headers sql.NullString, req_data_j string, req_data_t string, integration_config string, authors_selector string, summary_selector string, pub_date_selector string, title_selector string, url_selector string, stories_selector string, web_url string, finished chan bool) {
	// 	// CrawlUrl(rss_url)
	// 	// fmt.Println("Valid", rows)
	// 	fmt.Println("######last run######", last_run)
	// 	fmt.Println("####id########", id)
	// 	fmt.Println("#####rss url#######", rss_url)
	// 	fmt.Println("######status######", status)
	// 	fmt.Println("####freq########", category)
	// 	fmt.Println("####day range########", archive_day_range)
	// 	fmt.Println("####proxy ########", proxy_type)
	// 	fmt.Println("#####req method#######", req_method)
	// 	fmt.Println("#####req headers#######", req_headers)
	// 	fmt.Println("####req data########", req_data_j)
	// 	fmt.Println("#####req data t#######", req_data_t)
	// 	fmt.Println("######integration######", integration_config)
	// 	fmt.Println("######author######", authors_selector)
	// 	fmt.Println("#####summary#######", summary_selector)
	// 	fmt.Println("#####pub_date#######", pub_date_selector)
	// 	fmt.Println("######title######", title_selector)
	// 	fmt.Println("#####url#######", url_selector)
	// 	fmt.Println("#####stories#######", stories_selector)
	// 	fmt.Println("#####web url#######", web_url)

	// 	finished <- true
	// }(last_run, id, rss_url, status, category, archive_day_range, proxy_type, req_method, req_headers, req_data_j, req_data_t, integration_config, authors_selector, summary_selector, pub_date_selector, title_selector, url_selector, stories_selector, web_url, done)
	// }
	// for i := 0; i < 1; i++ {
	// 	fmt.Println(<-done)
	// }
}

func ExecuteWithoutGoroutineQuery(db *sql.DB, rawQuery string) {
	// done := make(chan bool)
	fmt.Println("***********")
	rows, _ := db.Query(rawQuery)
	db.Close()
	for rows.Next() {
		// rows.Scan(&authors_selector, &req_data_j, &req_data_t, &req_data_j, &integration_config, &authors_selector, &summary_selector)
		// fmt.Println(authors_selector, req_data_j, req_data_t, integration_config, authors_selector, summary_selector)
		rows.Scan(&authors_selector, &summary_selector, &req_data_j, &req_data_t)
		fmt.Println(authors_selector, summary_selector, req_data_j, req_data_t)
	}
}
