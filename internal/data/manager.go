package data

import (
	"log"
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
			Tags:        make(map[string][]string),
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

// GetPostsByTag 根据标签获取文章
func (m *Manager) GetPostsByTag(tag string) []*Post {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var posts []*Post
	seen := make(map[string]bool) // 用于去重

	if postIDs, exists := m.data.Tags[tag]; exists {
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
	}

	// 按创建时间倒序排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateTime.After(posts[j].CreateTime)
	})

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

// GetAllTags 获取所有标签
func (m *Manager) GetAllTags() map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tags := make(map[string]int)
	for tag, postIDs := range m.data.Tags {
		tags[tag] = len(postIDs)
	}
	return tags
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
func (m *Manager) UpdatePost(post *Post) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 更新文章
	m.data.Posts[post.ID] = post

	// 更新分类索引 - 避免重复添加
	if post.Category != "" {
		categoryPosts := m.data.Categories[post.Category]
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

	// 更新标签索引 - 避免重复添加
	for _, tag := range post.Tags {
		tagPosts := m.data.Tags[tag]
		found := false
		for _, id := range tagPosts {
			if id == post.ID {
				found = true
				break
			}
		}
		if !found {
			m.data.Tags[tag] = append(tagPosts, post.ID)
		}
	}

	m.data.LastUpdate = time.Now()
}

// RemovePost 移除文章
func (m *Manager) RemovePost(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if post, exists := m.data.Posts[id]; exists {
		// 从分类索引中移除
		if post.Category != "" {
			if postIDs, ok := m.data.Categories[post.Category]; ok {
				for i, postID := range postIDs {
					if postID == id {
						m.data.Categories[post.Category] = append(postIDs[:i], postIDs[i+1:]...)
						break
					}
				}
			}
		}

		// 从标签索引中移除
		for _, tag := range post.Tags {
			if postIDs, ok := m.data.Tags[tag]; ok {
				for i, postID := range postIDs {
					if postID == id {
						m.data.Tags[tag] = append(postIDs[:i], postIDs[i+1:]...)
						break
					}
				}
			}
		}

		// 从文章列表中移除
		delete(m.data.Posts, id)
		m.data.LastUpdate = time.Now()
	}
}

// Clear 清空所有数据
func (m *Manager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data.Posts = make(map[string]*Post)
	m.data.Categories = make(map[string][]string)
	m.data.Tags = make(map[string][]string)
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

	// 清理标签索引中的重复项
	for tag, postIDs := range m.data.Tags {
		seen := make(map[string]bool)
		var cleanedIDs []string
		for _, id := range postIDs {
			if !seen[id] {
				seen[id] = true
				cleanedIDs = append(cleanedIDs, id)
			}
		}
		m.data.Tags[tag] = cleanedIDs
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

// GetAllPosts 获取所有文章，按更新时间倒序排序
func (m *Manager) GetAllPosts() []*Post {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var posts []*Post
	for _, post := range m.data.Posts {
		posts = append(posts, post)
	}

	// 按更新时间倒序排序（时间最新的在前面）
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].UpdateTime.After(posts[j].UpdateTime)
	})

	return posts
}
