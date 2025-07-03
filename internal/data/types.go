package data

import (
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// CustomTime 自定义时间类型，支持多种格式
type CustomTime struct {
	time.Time
}

// UnmarshalYAML 自定义YAML解析
func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error {
	timeStr := strings.TrimSpace(value.Value)
	if timeStr == "" {
		ct.Time = time.Time{}
		return nil
	}

	// 支持的时间格式
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // 我们的格式
		"2006-01-02 15:04:05",       // 标准格式
		"2006-01-02",                // 仅日期
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			ct.Time = t
			return nil
		}
	}

	// 如果都解析失败，使用当前时间
	ct.Time = time.Now()
	return nil
}

// Post 文章数据结构
type Post struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	GitHubURL  string    `json:"github_url"`
	Content    string    `json:"content"`
	HTML       string    `json:"html"`
	Summary    string    `json:"summary"`
	Category   string    `json:"category"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	FilePath   string    `json:"file_path"`
}

// BlogData 博客数据结构
type BlogData struct {
	Posts       map[string]*Post    `json:"posts"`
	Categories  map[string][]string `json:"categories"`
	SearchIndex map[string][]string `json:"search_index"`
	LastUpdate  time.Time           `json:"last_update"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Posts    []*Post `json:"posts"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Query    string  `json:"query"`
}

// FrontMatter Front Matter 数据结构
type FrontMatter struct {
	Title       string     `yaml:"title"`
	Author      string     `yaml:"author"`
	GitHubURL   string     `yaml:"github_url"`
	CreatedAt   CustomTime `yaml:"created_at"`
	UpdatedAt   CustomTime `yaml:"updated_at"`
	Description string     `yaml:"description"`
}
