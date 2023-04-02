package constant

import "fmt"

const Name = "manoj"

const UrlPrifix = "https://www.manoj.com"

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "django_app"
	PASSWORD = "django"
	DBNAME   = "contify_db"

	SNAPSHOTDBNAME   = "rssfeed_snapshot"
	SNAPSHOTUSER     = "django_app"
	SNAPSHOTPASSWORD = "django"
)

var DbConfig = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DBNAME)

var SnapshotDbConfig = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, SNAPSHOTUSER, SNAPSHOTPASSWORD, SNAPSHOTDBNAME)
