package crawl

import (
	dbConnecter "awesomeGo/storage"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "awesomeGo/logging"

	"github.com/gosimple/slug"
	dps "github.com/markusmobius/go-dateparser"
)

var logger = log.New()

// add your db pipeline
type SnapshotData struct {
	Author       []string `json:"authors"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	PubDate      string   `json:"pub_date"`
	CompanyRss   string   `json:"company_rss"`
	BodyHtmlData string   `json:"body_html"`
}

func (a SnapshotData) Value() string {
	dt, _ := json.Marshal(a)
	return string(dt)
}

func addOrUpdateRecord(data []map[string]string, statusCode int, col *Column) {
	now := time.Now()
	db := dbConnecter.GetPostgresSnapshotDbConn()
	for count, row := range data {
		status := 1
		sqlStr := "INSERT INTO publications_rsssnapshot(url, rss_feed_id, data, created_on, updated_on, enrich_status, translation_status, status) VALUES "
		storyUrl := row["url"]
		title := row["title"]
		pubDate := row["pubDate"]
		if pubDate != "" {
			// parsedPubDate, _ := dps.Parse(nil, pubDate, "02 January 2006")
			// archiveDayRange := now.AddDate(0, 0, -col.archive_day_range)
			parsedPubDate, _ := dps.Parse(nil, pubDate, "02 January 2006")
			archiveDayRange := parsedPubDate.Time.AddDate(0, 0, col.archive_day_range)
			if archiveDayRange.After(now) {
				fmt.Println("date1 is before date2: ", archiveDayRange, now)
			}
			if archiveDayRange.Before(now) {
				logger.Info().Println("Date older then archive day range: ", pubDate, parsedPubDate, archiveDayRange, col.archive_day_range)
				status = 3
			}
		}
		summary := row["summary"]
		authorData := strings.Split(row["author"], ",")
		companyRss := row["company_rss"]
		slugText := getSlugFromText(title)
		attrs := new(SnapshotData)
		attrs.Author = authorData
		attrs.CompanyRss = companyRss
		attrs.PubDate = pubDate
		attrs.Slug = slugText
		attrs.Title = strings.TrimSpace(title)
		attrs.BodyHtmlData = strings.TrimSpace(summary)
		data := attrs.Value()
		sqlStr += fmt.Sprintf("('%v', %v, '%v', NOW(), NOW(), 1, 1, '%v')", storyUrl, col.id, data, status)

		res, err := db.Exec(sqlStr)
		if err != nil {
			logger.Info().Println("Got Error: ", err)
		}
		logger.Info().Println("###########", res, count)
	}

}

func getSlugFromText(text string) string {
	slugText := slug.Make(text)
	return slugText
}
