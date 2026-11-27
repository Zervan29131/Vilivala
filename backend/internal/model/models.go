package model

import (
	"time"
	"unicode"

	"gorm.io/gorm"
)

// User represents the blog administrator
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string    `gorm:"size:256;not null" json:"-"`
	Nickname     string    `gorm:"size:64" json:"nickname"`
	Avatar       string    `gorm:"size:512" json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
}

// Category represents a blog post category
type Category struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex:idx_cat_parent_name;size:64;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	ParentID  uint      `gorm:"uniqueIndex:idx_cat_parent_name;index;default:0" json:"parent_id"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Posts     []Post    `gorm:"foreignKey:CategoryID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// Tag represents a blog post tag
type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Posts     []Post    `gorm:"many2many:post_tags;" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// Post represents a blog post
type Post struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"size:256;not null" json:"title"`
	Slug       string    `gorm:"uniqueIndex;size:256;not null" json:"slug"`
	Summary    string    `gorm:"size:1024" json:"summary"`
	Content    string    `gorm:"type:longtext" json:"content,omitempty"`
	CoverImage string    `gorm:"size:512" json:"cover_image"`
	Status     string    `gorm:"size:16;default:draft;index" json:"status"` // draft | published
	ViewCount  uint      `gorm:"default:0" json:"view_count"`
	CategoryID *uint     `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"-"`
	// CategoryOut is the API representation of the category, with breadcrumb path (filled by handlers).
	CategoryOut *CategoryRef `gorm:"-" json:"category,omitempty"`
	IsTop       bool         `gorm:"default:false;index" json:"is_top"`
	WordCount   int          `gorm:"default:0" json:"word_count"`
	ReadingTime int          `gorm:"-" json:"reading_time"`
	Tags        []Tag        `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// PathNode is one step of a category breadcrumb.
type PathNode struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CategoryRef is the category shape embedded in a post.
type CategoryRef struct {
	ID   uint       `json:"id"`
	Name string     `json:"name"`
	Slug string     `json:"slug"`
	Path []PathNode `json:"path"`
}

// Page is a standalone Markdown page (about, links, ...).
type Page struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug      string    `gorm:"uniqueIndex;size:64;not null" json:"slug"`
	Title     string    `gorm:"size:64;not null" json:"title"`
	Content   string    `gorm:"type:longtext" json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SiteSetting stores singleton site-wide dynamic configuration as JSON (ID = 1).
type SiteSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Payload   string    `gorm:"type:longtext" json:"payload"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AfterFind derives reading_time (minutes) from the stored word count.
func (p *Post) AfterFind(_ *gorm.DB) error {
	p.ReadingTime = ReadingMinutes(p.WordCount)
	return nil
}

// ReadingMinutes = max(1, ceil(words/400)).
func ReadingMinutes(words int) int {
	m := (words + 399) / 400
	if m < 1 {
		m = 1
	}
	return m
}

// CountWords counts non-whitespace characters.
func CountWords(s string) int {
	n := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			n++
		}
	}
	return n
}
