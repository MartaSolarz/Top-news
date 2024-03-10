package newsDB

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"top-news/backend/internal/adapter"
	"top-news/backend/internal/models"
)

type DBOperations struct {
	mysqlConn *adapter.DBConnection
	tableName string
	TTL       int
}

func NewDBOperations(mysqlConn *adapter.DBConnection, tableName string, ttl int) *DBOperations {
	return &DBOperations{
		mysqlConn: mysqlConn,
		tableName: tableName,
		TTL:       ttl,
	}
}

func (dbOps *DBOperations) GetNews() ([]*models.Article, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT * FROM %s", dbOps.tableName))
	if err != nil {

	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("could not close rows: %s", err.Error())
		}
	}(rows)

	articles := make([]*models.Article, 0)
	for rows.Next() {
		var id, thbWidth, thbHeight int
		var website, copyRight, title, description, summary, authors, sourceUrl, thbUrl string
		var publishDateStr, expireAtStr []byte
		if err = rows.Scan(
			&id, &website, &copyRight, &title, &description, &summary,
			&publishDateStr, &sourceUrl, &authors, &thbUrl, &thbWidth, &thbHeight, &expireAtStr,
		); err != nil {
			log.Printf("could not scan row: %s", err.Error())
			continue
		}

		publishDate, err := time.Parse("2006-01-02 15:04:05", string(publishDateStr))
		if err != nil {
			log.Printf("could not parse publish date: %s", err.Error())
			continue
		}
		expireAt, err := time.Parse("2006-01-02 15:04:05", string(expireAtStr))
		if err != nil {
			log.Printf("could not parse expire date: %s", err.Error())
			continue
		}
		record := models.NewArticleFromDB(
			id, thbWidth, thbHeight, website, copyRight, title, description, summary,
			sourceUrl, authors, thbUrl, publishDate, expireAt,
		)

		articles = append(articles, record)
	}
	return articles, nil
}

func (dbOps *DBOperations) GetTitles() ([]string, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT title FROM %s", dbOps.tableName))
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("could not close rows: %s", err.Error())
		}
	}(rows)

	titles := make([]string, 0)
	for rows.Next() {
		var title string
		if err = rows.Scan(&title); err != nil {
			log.Printf("could not scan row: %s", err.Error())
			continue
		}
		titles = append(titles, title)
	}
	return titles, nil
}

func (dbOps *DBOperations) PutNews(articles []*models.Article) error {
	query := fmt.Sprintf("INSERT INTO %s (website, copy_right, title, description, summary, publish_date, source_url, authors, thumbnail_url, thumbnail_width, thumbnail_height, expire_date) VALUES ", dbOps.tableName)
	var params []interface{}
	var placeholders []string

	for _, article := range articles {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, article.Website, article.CopyRight, article.Title, article.Description, article.Summary, article.PublishDate, article.SourceURL, article.Authors, article.Thumbnail.URL, article.Thumbnail.Width, article.Thumbnail.Height, article.ExpireAt)
	}

	query += strings.Join(placeholders, ", ")

	stmt, err := dbOps.mysqlConn.Client.Prepare(query)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Printf("could not close statement: %s", err.Error())
			return
		}
	}()

	_, err = stmt.Exec(params...)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	log.Printf("successfully inserted %d records", len(articles))

	return nil
}

func (dbOps *DBOperations) GetTTL() int {
	return dbOps.TTL
}
