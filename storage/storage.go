package storage

import (
	constant "awesomeGo/constant"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	url                string
	last_run           string
	data               string
	rss_url            string
	stories_selector   string
	url_selector       string
	title_selector     string
	pub_date_selector  string
	summary_selector   string
	authors_selector   string
	integration_config string
	req_data_j         string
	req_data_t         string
	req_headers        string
	id, status         int64
	req_method         int
	proxy_type         int
	archive_day_range  int
	frequency          int
)

func GetPostgresDbConn() *sql.DB {
	connStr := constant.DbConfig
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if db.Ping() != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

func GetPostgresSnapshotDbConn() *sql.DB {
	connStr := constant.SnapshotDbConfig
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("connStr", connStr)
		fmt.Println(err)
		panic(err)
	}
	if db.Ping() != nil {
		fmt.Println("connStr2", connStr)
		fmt.Println(err)
		panic(err)
	}

	return db
}

func ExecuteQuery(db *sql.DB, rawQuery string) {
	done := make(chan bool)
	rows, _ := db.Query(rawQuery)
	db.Close()
	for rows.Next() {
		rows.Scan(&last_run, &id, &url, &status, &data, &frequency, &archive_day_range, &proxy_type, &req_method,
			&req_headers, &req_data_j, &req_data_t, &req_data_j, &integration_config, &authors_selector,
			&summary_selector, &pub_date_selector, &title_selector, &url_selector, &stories_selector, &rss_url, &data)
		go func(id int64, finished chan bool) {
			fmt.Println("before wait: ", id)
			// time.Sleep(1 * time.Second)
			fmt.Println("after wait: ", id)
			finished <- true
		}(id, done)
	}
	for i := 0; i < 5; i++ {
		fmt.Println(<-done)
	}

}
