package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Manager 数据管理器
type Manager struct {
	data     *BlogData
	postsDir string
	mutex    sync.RWMutex
}

// NewManager 创建新的数据管理器
func NewManager(postsDir string) *Manager {
	return &Manager{
		data: &BlogData{
			Posts:       make(map[string]*Post),
			Categories:  make(map[string][]string),
			SearchIndex: make(map[string][]string),
			LastUpdate:  time.Now(),
		},
		postsDir: postsDir,
	}
}

// GetData 获取博客数据
func (m *Manager) GetData() *BlogData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data
}

// GetPostsDir 获取 posts 目录路径
func (m *Manager) GetPostsDir() string {
	return m.postsDir
}

// GetPost 根据ID获取文章
func (m *Manager) GetPost(id string) (*Post, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	post, exists := m.data.Posts[id]
	return post, exists
}

// GetPostsByCategory 根据分类获取文章
func (m *Manager) GetPostsByCategory(category string) []*Post {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var posts []*Post
	seen := make(map[string]bool) // 用于去重

	log.Printf("查询分类 %s 的文章", category)

	if postIDs, exists := m.data.Categories[category]; exists {
		log.Printf("找到分类 %s，包含 %d 个文章ID: %v", category, len(postIDs), postIDs)
		for _, id := range postIDs {
			// 跳过重复的ID
			if seen[id] {
				continue
			}
			seen[id] = true

			if post, ok := m.data.Posts[id]; ok {
				// 跳过about文章，它只在About页面显示
				if post.ID != "about" {
					posts = append(posts, post)
				}
			}
		}
	} else {
		log.Printf("未找到分类 %s", category)
		// 打印所有可用的分类
		log.Printf("可用分类: %v", func() []string {
			var cats []string
			for cat := range m.data.Categories {
				cats = append(cats, cat)
			}
			return cats
		}())
	}

	// 按创建时间倒序排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateTime.After(posts[j].CreateTime)
	})

	log.Printf("最终返回 %d 篇文章", len(posts))
	return posts
}

// GetAllCategories 获取所有分类
func (m *Manager) GetAllCategories() map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	categories := make(map[string]int)
	for category, postIDs := range m.data.Categories {
		categories[category] = len(postIDs)
	}
	return categories
}

// Search 搜索文章
func (m *Manager) Search(query string, page, pageSize int) *SearchResult {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var results []*Post
	query = strings.ToLower(query)

	// 简单的模糊搜索实现，排除about文章
	for _, post := range m.data.Posts {
		// 跳过about文章，它只在About页面显示
		if post.ID == "about" {
			continue
		}
		if strings.Contains(strings.ToLower(post.Title), query) ||
			strings.Contains(strings.ToLower(post.Content), query) {
			results = append(results, post)
		}
	}

	// 按更新时间倒序排序（时间最新的在前面）
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdateTime.After(results[j].UpdateTime)
	})

	total := len(results)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return &SearchResult{
			Posts:    []*Post{},
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Query:    query,
		}
	}

	if end > total {
		end = total
	}

	return &SearchResult{
		Posts:    results[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Query:    query,
	}
}

// UpdatePost 更新文章
func (m *Manager) UpdatePost(post *Post) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 删除旧的分类映射
	if existingPost, exists := m.data.Posts[post.ID]; exists {
		oldCategory := existingPost.Category
		if postIDs, ok := m.data.Categories[oldCategory]; ok {
			for i, id := range postIDs {
				if id == post.ID {
					m.data.Categories[oldCategory] = append(postIDs[:i], postIDs[i+1:]...)
					break
				}
			}
		}
	}

	// 更新分类映射
	if post.Category != "" {
		categoryPosts := m.data.Categories[post.Category]
		// 检查是否已存在，避免重复添加
		found := false
		for _, id := range categoryPosts {
			if id == post.ID {
				found = true
				break
			}
		}
		if !found {
			m.data.Categories[post.Category] = append(categoryPosts, post.ID)
		}
	}

	// 更新文章
	m.data.Posts[post.ID] = post
	m.data.LastUpdate = time.Now()

	return nil
}

// DeletePost 删除文章
func (m *Manager) DeletePost(postID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	post, exists := m.data.Posts[postID]
	if !exists {
		return fmt.Errorf("文章不存在: %s", postID)
	}

	// 从分类中删除
	if post.Category != "" {
		if postIDs, ok := m.data.Categories[post.Category]; ok {
			for i, id := range postIDs {
				if id == postID {
					m.data.Categories[post.Category] = append(postIDs[:i], postIDs[i+1:]...)
					break
				}
			}
		}
	}

	// 删除文章
	delete(m.data.Posts, postID)
	m.data.LastUpdate = time.Now()

	return nil
}

// cleanup 清理无效的索引
func (m *Manager) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data.Categories = make(map[string][]string)

	// 重新构建索引
	for _, post := range m.data.Posts {
		// 重新构建分类索引
		if post.Category != "" {
			m.data.Categories[post.Category] = append(m.data.Categories[post.Category], post.ID)
		}
	}

	// 清理空的索引
	for category, postIDs := range m.data.Categories {
		var cleanedIDs []string
		for _, id := range postIDs {
			if _, exists := m.data.Posts[id]; exists {
				cleanedIDs = append(cleanedIDs, id)
			}
		}
		if len(cleanedIDs) == 0 {
			delete(m.data.Categories, category)
		} else {
			m.data.Categories[category] = cleanedIDs
		}
	}
}

// Clear 清空所有数据
func (m *Manager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data.Posts = make(map[string]*Post)
	m.data.Categories = make(map[string][]string)
	m.data.SearchIndex = make(map[string][]string)
	m.data.LastUpdate = time.Now()
}

// CleanupDuplicates 清理重复的文章ID
func (m *Manager) CleanupDuplicates() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 清理分类索引中的重复项
	for category, postIDs := range m.data.Categories {
		seen := make(map[string]bool)
		var cleanedIDs []string
		for _, id := range postIDs {
			if !seen[id] {
				seen[id] = true
				cleanedIDs = append(cleanedIDs, id)
			}
		}
		m.data.Categories[category] = cleanedIDs
	}

	m.data.LastUpdate = time.Now()
}

// GenerateID 根据文件路径生成文章ID
func (m *Manager) GenerateID(filePath string) string {
	// 移除扩展名和posts目录前缀
	relPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	if strings.HasPrefix(relPath, m.postsDir+"/") {
		relPath = strings.TrimPrefix(relPath, m.postsDir+"/")
	}
	return relPath
}

// GetCategoryFromPath 从文件路径获取分类
func (m *Manager) GetCategoryFromPath(filePath string) string {
	relPath := filePath
	if strings.HasPrefix(relPath, m.postsDir+"/") {
		relPath = strings.TrimPrefix(relPath, m.postsDir+"/")
	}

	parts := strings.Split(relPath, "/")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

// GetAllPosts 获取所有文章（按更新时间倒序排序）
func (m *Manager) GetAllPosts() []*Post {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var posts []*Post
	for _, post := range m.data.Posts {
		// 跳过about文章，它只在About页面显示
		if post.ID != "about" {
			posts = append(posts, post)
		}
	}

	// 按更新时间倒序排序（时间最新的在前面）
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].UpdateTime.After(posts[j].UpdateTime)
	})

	return posts
}

// InitializeData 初始化数据，遍历并加载所有Markdown文件
func (m *Manager) InitializeData(parser Parser) error {
	// 遍历posts目录
	err := filepath.Walk(m.postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理.md文件
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// 检查是否是posts根目录下的文件
		relPath, err := filepath.Rel(m.postsDir, path)
		if err != nil {
			log.Printf("计算相对路径失败 %s: %v", path, err)
			return nil
		}

		// 如果是posts根目录下的文件（不包含路径分隔符），只允许about.md
		if !strings.Contains(relPath, string(filepath.Separator)) {
			filename := filepath.Base(path)
			if filename != "about.md" {
				log.Printf("跳过posts根目录下的文件: %s", path)
				return nil
			}
		}

		// 解析Markdown文件
		post, err := parser.ParseFile(path, m.postsDir)
		if err != nil {
			log.Printf("解析文件失败 %s: %v", path, err)
			return nil // 继续处理其他文件
		}

		// 添加到数据管理器
		m.UpdatePost(post)
		log.Printf("加载文章: %s", post.Title)

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历posts目录失败: %w", err)
	}

	// 清理重复数据
	m.CleanupDuplicates()

	log.Printf("数据初始化完成，共加载 %d 篇文章", len(m.data.Posts))
	return nil
}

// Parser 定义Markdown解析器接口
type Parser interface {
	ParseFile(filePath, postsDir string) (*Post, error)
}
