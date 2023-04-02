package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Snapshotlog struct {
	url    string `json:"URL"`
	id     int64  `json:"ID"`
	status int64  `json:"Status"`
	data   string `json:"Data"`
}

var (
	urllog           string
	idlog, statuslog int64
	datalog          string
	snaplog          []Snapshotlog
)

func GetDbConnlog() {
	fmt.Println("Getting db connection")
	dsn := "host=localhost user=postgres password=postgres dbname=rssfeed_snapshot port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var result int64
	var s int64
	db.Table("publications_rsssnapshot").Where("status = ?", 1).Count(&result)
	// return result

	db.Raw("select count(*) from publications_rsssnapshot where status=1").Scan(&s)

	fmt.Println("Connect Count: ", result)
	fmt.Println("Connect Count Rows: ", s)
}

func displaylog(rows *sql.Rows) {
	for rows.Next() {
		rows.Scan(&id, &url, &status, &data)
		fmt.Println(id)
		snaplog = append(snaplog, Snapshotlog{url: url, id: id, status: status, data: data})
	}
}

func GetPostgresDbConnlog() *sql.DB {
	// connStr := constant.DbConfig
	db, err := sql.Open("postgres", "host=localhost user=postgres password=postgres dbname=rssfeed_snapshot port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	if db.Ping() != nil {
		panic(err)
	}

	return db
}

func ExecuteQuerylog(db *sql.DB, rawQuery string) {
	rows, _ := db.Query(rawQuery)
	db.Close()
	for rows.Next() {
		rows.Scan(&id, &url, &status, &data)
		go func(id int64) {
			fmt.Println("before wait: ", id)
			time.Sleep(1 * time.Second)
			fmt.Println("after wait: ", id)
		}(id)
	}
	// 	// 	// snap = append(snap, Snapshot{url:url, id:id, status:status, data:data})
	// }

	// go display(rows)
	time.Sleep(2 * time.Second)
	// fmt.Println("Snapshot: ", snap)

}

func CrawlUrl() {
	urlCh := make(chan int)
	var myUrlArray [1]string
	// append(myUrlArray, "https://www.europapress.es/castilla-lamancha/noticia-lm-priorizara-empresas-incluyan-menus-productos-ecologicos-centros-sociosanitarios-20230214131303.html")
	// urls := [2]string{"https://www.europapress.es/castilla-lamancha/noticia-lm-priorizara-empresas-incluyan-menus-productos-ecologicos-centros-sociosanitarios-20230214131303.html", "https://www.newsnow.co.uk/h/Technology"}
	for _, url := range myUrlArray {
		go func(url string) {

			c := colly.NewCollector()
			c.OnRequest(func(r *colly.Request) {
				r.Headers.Set("Accept-Language", "en-US;q=0.9")
				fmt.Println("Visiting", r.URL)
			})
			c.OnResponse(func(r *colly.Response) {
				fmt.Println("Status:", r.StatusCode)
				// fmt.Printf("Type Status: %T", r.StatusCode)
				urlCh <- r.StatusCode
			})
			c.Visit(url)
		}(url)
	}
	//fmt.Println("I am from utils", c.Name)
	for range myUrlArray {
		fmt.Println("I am from urls loop")
		fmt.Println(<-urlCh)
	}
	fmt.Println("######last run######", last_run)
	fmt.Println("####id########", id)
	fmt.Println("#####rss url#######", rss_url)
	fmt.Println("######status######", status)
	//fmt.Println("####freq########", category)
	fmt.Println("####day range########", archive_day_range)
	fmt.Println("####proxy ########", proxy_type)
	fmt.Println("#####req method#######", req_method)
	fmt.Println("#####req headers#######", req_headers)
	fmt.Println("####req data########", req_data_j)
	fmt.Println("#####req data t#######", req_data_t)
	fmt.Println("######integration######", integration_config)
	fmt.Println("######author######", authors_selector)
	fmt.Println("#####summary#######", summary_selector)
	fmt.Println("#####pub_date#######", pub_date_selector)
	fmt.Println("######title######", title_selector)
	fmt.Println("#####url#######", url_selector)
	fmt.Println("#####stories#######", stories_selector)
	//fmt.Println("#####web url#######", web_url)
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
