package main

import (
	// util "awesomeGo/utils"

	crawl "awesomeGo/crawl2"
	dbConnecter "awesomeGo/storage"
	"flag"
	"fmt"
	"time"
)

// util "awesomeGo/utils"

// func main() {
// 	appLogger := New()

func main() {
	start := time.Now()
	fmt.Println("#####################Started Main####################", start)

	category := flag.Int("category", 10, "category option passed from cmd")
	integrtion_type := flag.String("integrtion_type", "3", "category option passed from cmd")
	flag.Parse()
	sinceTime := time.Now().Add(time.Duration(-*category) * time.Minute).UTC().Format("2006-01-02T15:04:05-0700")
	rawQuery := fmt.Sprintf(`SELECT id, rss_url,stories_selector,url_selector,title_selector,pub_date_selector,summary_selector,
	authors_selector,archive_day_range, req_headers
	FROM publications_rssfeed
	WHERE(
		active = True AND integration_type IN (%v) AND show_for_client_id = 405506
		AND (last_run <= '%v' OR last_run IS NULL)
	) ORDER BY last_run ASC limit 10;`, *integrtion_type, sinceTime)

	// rawQuery1 := "select id, rss_url,stories_selector,url_selector,title_selector,pub_date_selector,summary_selector,authors_selector,proxy_type, req_headers FROM publications_rssfeed where  integration_type=3 and active=true and show_for_client_id=405506 and req_headers is not null limit 1"
	db := dbConnecter.GetPostgresDbConn()
	dbRows := crawl.ExecuteQuery(db, rawQuery)
	crawl.ProcessData(dbRows)

	fmt.Println("#####################Fineshed Main###################", time.Since(start))
}
