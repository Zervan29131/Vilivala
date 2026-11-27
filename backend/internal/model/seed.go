package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed inserts default admin and demo data when the users table is empty.
func Seed(db *gorm.DB) error {
	var n int64
	if err := db.Model(&User{}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	log.Println("[Seed] users table empty, seeding default data")

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := db.Create(&User{Username: "admin", PasswordHash: string(hash), Nickname: "博主"}).Error; err != nil {
		return err
	}

	cats := []Category{{Name: "技术笔记", Slug: "tech-notes"}, {Name: "生活随笔", Slug: "life-essay"}}
	tags := []Tag{{Name: "Go", Slug: "go"}, {Name: "Vue", Slug: "vue"}, {Name: "Markdown", Slug: "markdown"}}
	if err := db.Create(&cats).Error; err != nil {
		return err
	}
	if err := db.Create(&tags).Error; err != nil {
		return err
	}

	p1 := Post{
		Title:      "Go 语言入门指南",
		Slug:       "go-getting-started",
		Summary:    "从零开始认识 Go：安装、第一个程序与并发初体验。",
		Status:     "published",
		CategoryID: &cats[0].ID,
		Tags:       []Tag{tags[0]},
		Content:    "# Go 语言入门\n\n## 安装\n\n```bash\nbrew install go\ngo version\n```\n\n## Hello World\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, Go!\")\n}\n```\n\n## 特点\n\n- 编译快\n- 原生并发（goroutine）\n- 标准库丰富\n\n> Go 的设计哲学：简单、直接、够用。\n\n| 特性 | 说明 |\n| --- | --- |\n| goroutine | 轻量级线程 |\n| channel | 协程间通信 |\n",
	}
	p2 := Post{
		Title:      "用 Vue 3 和 Markdown 写博客",
		Slug:       "vue3-markdown-blog",
		Summary:    "记录搭建个人博客的思路：Vue 3 + Go + MySQL。",
		Status:     "published",
		CategoryID: &cats[1].ID,
		Tags:       []Tag{tags[1], tags[2]},
		Content:    "# 搭建个人博客\n\n## 技术栈\n\n1. Vue 3 + TypeScript\n2. Go + Gin + GORM\n3. MySQL 8.0\n\n## 示例代码\n\n```ts\nconst count = ref(0)\nfunction inc() {\n  count.value++\n}\n```\n\n## 小结\n\n接口契约先行，前后端联调会顺利很多。\n",
	}
	return db.Create(&[]Post{p1, p2}).Error
}
