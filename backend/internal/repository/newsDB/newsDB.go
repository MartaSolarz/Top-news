package newsDB

import (
	"database/sql"
	"fmt"
	"log"
	"top-news/backend/internal/adapter"
	"top-news/backend/internal/models"
)

type DBOperations struct {
	mysqlConn *adapter.DBConnection
	tableName string
}

func NewDBOperations(mysqlConn *adapter.DBConnection, tableName string) *DBOperations {
	return &DBOperations{
		mysqlConn: mysqlConn,
		tableName: tableName,
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
		var id int
		var website, copyRight, title, description, summary, publishDate, sourceUrl, thbUrl, thbWidth, thbHeight string
		if err = rows.Scan(
			&id, &website, &copyRight, &title, &description, &summary,
			&publishDate, &sourceUrl, &thbUrl, &thbWidth, &thbHeight,
		); err != nil {
			log.Printf("could not scan row: %s", err.Error())
			continue
		}
		record := models.NewArticleFromDB(
			id, website, copyRight, title, description, summary,
			publishDate, sourceUrl, thbUrl, thbWidth, thbHeight,
		)

		articles = append(articles, record)
	}
	return articles, nil
}

func (dbOps *DBOperations) PutNews(articles []*models.Article) error {
	query := fmt.Sprintf("INSERT INTO %s (website, copy_right, title, description, summary, publish_date, source_url, thumbnail_url, thumbnail_width, thumbnail_height) VALUES ", dbOps.tableName)
	for i, article := range articles {
		query += fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			article.Website, article.CopyRight, article.Title, article.Description, article.Summary,
			article.PublishDate, article.SourceURL, article.Thumbnail.URL, article.Thumbnail.Width, article.Thumbnail.Height,
		)
		if i != len(articles)-1 {
			query += ", "
		}
	}
	_, err := dbOps.mysqlConn.Client.Exec(query)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return nil
}
