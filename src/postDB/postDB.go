package postDB

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var PageSize int = 3
var host string = "posts/post"

type Post struct {
	ID      int64     `bun:",pk,autoincrement"`
	Link    string    `bun:"link"`
	Title   string    `bun:"title"`
	Preview string    `bun:"preview"`
	Content string    `bun:"content"`
	PubDate time.Time `bun:"timestamp"`
}

func ConnectDB() *bun.DB {
	sqldb, err := sql.Open("postgres", "user=postgres password=1 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return bun.NewDB(sqldb, pgdialect.New())
}

func InsertPost(db *bun.DB, post *Post) {
	id, _ := CountPosts(db)
	id++
	post.Link = fmt.Sprintf("%s?id=%d", host, id)
	if len(post.Content) > 50 {
		post.Preview = fmt.Sprintf("%s...", post.Content[:50])
	} else {
		post.Preview = post.Content
	}
	post.PubDate = time.Now()
	_, err := db.NewInsert().Model(post).Exec(context.Background())
	if err != nil {
		log.Printf("Error inserting post to DB: %v", err)
	}
}

func GetPosts(db *bun.DB, page int) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * PageSize

	err := db.NewSelect().
		Model(&posts).
		Order("id ASC").
		Limit(PageSize).
		Offset(offset).
		Scan(context.Background())

	return posts, err
}

func CountPosts(db *bun.DB) (int64, error) {
	var count int64
	err := db.NewSelect().
		Model((*Post)(nil)).
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &count)
	return count, err
}

func GetOnePost(db *bun.DB, id int) (Post, error) {
	var post Post
	err := db.NewSelect().
		Model(&post).
		Where("id = ?", id).
		Scan(context.Background())

	return post, err
}
