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
	mysqlConn         *adapter.DBConnection
	articlesTableName string
	emailsTableName   string
	TTL               int
	myDomain          string
}

func NewDBOperations(mysqlConn *adapter.DBConnection, articlesTableName, emailsTableName, myDomain string, ttl int) *DBOperations {
	return &DBOperations{
		mysqlConn:         mysqlConn,
		articlesTableName: articlesTableName,
		emailsTableName:   emailsTableName,
		TTL:               ttl,
		myDomain:          myDomain,
	}
}

func (dbOps *DBOperations) GetNews() ([]*models.Article, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT * FROM %s", dbOps.articlesTableName))
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

func (dbOps *DBOperations) GetNewsFromToday() ([]*models.Article, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT title, description, publish_date, source_url FROM %s WHERE DATE(publish_date) = CURDATE()", dbOps.articlesTableName))
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("could not close rows: %s", err.Error())
		}
	}(rows)

	articles := make([]*models.Article, 0)
	for rows.Next() {
		var title, description, sourceUrl string
		var publishDate []byte
		if err = rows.Scan(
			&title, &description, &publishDate, &sourceUrl,
		); err != nil {
			log.Printf("could not scan row: %s", err.Error())
			continue
		}

		record := &models.Article{
			Title:       title,
			Description: description,
			PublishDate: string(publishDate),
			SourceURL:   sourceUrl,
			MyDomain:    dbOps.myDomain,
		}

		articles = append(articles, record)
	}

	log.Printf("found %d articles", len(articles))

	return articles, nil
}

func (dbOps *DBOperations) GetTitles() ([]string, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT title FROM %s", dbOps.articlesTableName))
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
	query := fmt.Sprintf("INSERT INTO %s (website, copy_right, title, description, summary, publish_date, source_url, authors, thumbnail_url, thumbnail_width, thumbnail_height, expire_date) VALUES ", dbOps.articlesTableName)
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

func (dbOps *DBOperations) GetFavorites(ids [][]byte) ([]*models.Article, error) {
	var placeholders []string
	var args []interface{}
	for _, id := range ids {
		placeholders = append(placeholders, "?")
		args = append(args, string(id))
	}
	placeholdersStr := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", dbOps.articlesTableName, placeholdersStr)
	rows, err := dbOps.mysqlConn.Client.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
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

func (dbOps *DBOperations) SaveEmail(email string) error {
	query := fmt.Sprintf("INSERT INTO %s (email, created_at) VALUES (?, NOW())", dbOps.emailsTableName)
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

	_, err = stmt.Exec(email)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}

func (dbOps *DBOperations) GetEmails() ([]string, error) {
	rows, err := dbOps.mysqlConn.Client.Query(fmt.Sprintf("SELECT email FROM %s", dbOps.emailsTableName))
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("could not close rows: %s", err.Error())
		}
	}(rows)

	emails := make([]string, 0)
	for rows.Next() {
		var email string
		if err = rows.Scan(&email); err != nil {
			log.Printf("could not scan row: %s", err.Error())
			continue
		}
		emails = append(emails, email)
	}
	return emails, nil
}
