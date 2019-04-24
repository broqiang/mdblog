package mdfile

import (
	"errors"
	"sort"
	"strings"
)

// ListMap 将博客保存在 map 中，
type ListMap struct {
	Articles   Articles
	Categories []Category
	Tags       Tags
}

// CategoriesAll 获取所有的分类列表
func (list *ListMap) CategoriesAll() Categories {

	return list.Categories
}

// TagsAll 获取所有标签
func (list *ListMap) TagsAll() Tags {
	return list.Tags
}

// ArticlesAll 获取所有文章的列表
func (list *ListMap) ArticlesAll() Articles {
	return list.Articles
}

// ArticleByPath 根据文章的 Path 查询指定的文章
func (list *ListMap) ArticleByPath(path string) (Article, error) {
	if path != "" {
		for _, article := range list.Articles {
			if article.Path == strings.Trim(path, "/") {
				return article, nil
			}
		}
	}

	return Article{}, errors.New("can not found article")
}

func newListMap() *ListMap {

	list := ListMap{
		Categories: parseCategories(),
		Tags:       make(Tags, 0),
	}

	list.initArticles()

	list.initTags()

	return &list
}

// 初始化文章
func (list *ListMap) initArticles() {
	articles := make(Articles, 0)

	for i := 0; i < len(list.Categories); i++ {
		a := getAritclesSpecifiedCategory(&list.Categories[i])
		mergeArticles := make(Articles, len(articles)+len(a))
		copy(mergeArticles, articles)
		copy(mergeArticles[len(articles):], a)
		articles = mergeArticles
	}

	// 接收到所有文章，并按照 UpdateAt 倒序保存
	sort.Sort(&articles)
	list.Articles = articles
}

// 初始化 tags， 顺便做下排序
func (list *ListMap) initTags() {
	articles := list.Articles
	tags := list.Tags

	for _, article := range articles {
		for _, title := range article.Tags {
			i := getTagByTitle(tags, title)
			if i >= 0 {
				tags[i].Number++
			} else {
				tags = append(tags, Tag{
					Title:  strings.ToLower(title),
					Number: 1,
				})
			}

		}
	}

	list.Tags = tags

}

// 根据标签的标题获取标签
func getTagByTitle(tags Tags, title string) int {
	for i := 0; i < len(tags); i++ {
		if tags[i].Title == strings.ToLower(title) {
			return i
		}
	}

	return -1
}
