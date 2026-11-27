package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
	"vilivala-backend/internal/router"
)

func main() {
	// 1. Load configuration from .env / OS environment
	config.Load()

	// 2. Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// 3. Initialize database connection
	config.InitDB()

	// 4. Run AutoMigrate (creates / updates tables)
	dropLegacyIndexes()
	if err := autoMigrate(); err != nil {
		log.Fatalf("[Migration] AutoMigrate failed: %v", err)
	}

	if err := model.Seed(config.DB); err != nil {
		log.Fatalf("[Seed] failed: %v", err)
	}

	// 5. Set up routes
	r := router.Setup()

	// 6. Start server
	addr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	log.Printf("[Server] Starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[Server] Failed to start: %v", err)
		os.Exit(1)
	}
}

// dropLegacyIndexes removes the old global-unique category name index
// (names are now unique per parent only).
func dropLegacyIndexes() {
	m := config.DB.Migrator()
	if m.HasTable(&model.Category{}) && m.HasIndex(&model.Category{}, "idx_categories_name") {
		if err := m.DropIndex(&model.Category{}, "idx_categories_name"); err != nil {
			log.Printf("[Migration] drop idx_categories_name: %v", err)
		}
	}
}

func autoMigrate() error {
	db := config.DB
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Post{},
		&model.Page{},
	); err != nil {
		return err
	}
	ensureFulltextIndex()
	return backfillWordCount()
}

// ensureFulltextIndex adds a MySQL ngram FULLTEXT index on posts(title, summary, content) if missing.
func ensureFulltextIndex() {
	var cnt int64
	config.DB.Raw(
		"SELECT COUNT(1) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'posts' AND index_name = 'idx_posts_fulltext'",
	).Scan(&cnt)
	if cnt == 0 {
		if err := config.DB.Exec(
			"ALTER TABLE posts ADD FULLTEXT INDEX idx_posts_fulltext (title, summary, content) WITH PARSER ngram",
		).Error; err != nil {
			log.Printf("[Migration] fulltext index skipped: %v", err)
		} else {
			log.Println("[Migration] created FULLTEXT ngram index idx_posts_fulltext")
		}
	}
}

// backfillWordCount fills word_count for posts created before the column existed.
func backfillWordCount() error {
	var posts []model.Post
	if err := config.DB.Select("id", "content").Where("word_count = 0 AND content <> ''").Find(&posts).Error; err != nil {
		return err
	}
	for _, p := range posts {
		config.DB.Model(&model.Post{}).Where("id = ?", p.ID).UpdateColumn("word_count", model.CountWords(p.Content))
	}
	return nil
}
