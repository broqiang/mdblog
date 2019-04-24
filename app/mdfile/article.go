package mdfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/broqiang/mdblog/app/mylog"

	"github.com/broqiang/mdblog/app/helper"
)

// Articles 文章列表
type Articles []Article

// Article 文章内容
type Article struct {
	// 文章的标题
	Title string

	// 作者姓名
	Author string

	// 创建时间
	CreatedAt time.Time `toml:"created_at"`

	// 最后更新时间
	UpdatedAt time.Time `toml:"updated_at"`

	// 标签
	Tags []string

	// 所属分类的名称
	Category string

	// 头部图片 URL 地址
	HeadImg string `toml:"head_img"`

	// 作者的个人主页
	HomePage string `toml:"home_page"`

	// 简短的描述
	Description string

	// 文章主题内容， markdown
	Body string

	// 文章在服务器上的文件路由
	Path string
}

// 获取指定分类下的文章
func getAritclesSpecifiedCategory(category *Category) Articles {
	// 获取分类下文件的数量
	path := filepath.Join(getRootPath(), category.Path)

	// 确认目录是否为空，并且目录下是否有文件
	filesInfo, err := ioutil.ReadDir(path)
	if os.IsNotExist(err) {
		return nil
	}

	helper.PanicErr(err)

	number := 0

	articles := make([]Article, 0)

	for _, info := range filesInfo {
		// 如果目录下的是文件，就不处理
		if info.IsDir() {
			continue
		}

		// 获取文件名
		fileName := info.Name()
		// 获取文件名的后缀
		ext := filepath.Ext(fileName)

		// 如果后缀名不是 .md 就不处理
		if ext != ".md" {
			continue
		}

		// 获取文章 struct ， 如果是空的就不处理，继续下一次
		article := getArticleContent(filepath.Join(path, fileName))
		if reflect.DeepEqual(article, Article{}) {
			continue
		}

		// 将分类添加到文章中
		article.Category = category.Title
		article.Path = filepath.Join("/posts", strings.TrimSuffix(info.Name(), ext))

		articles = append(articles, article)

		// 文章数量 +1
		number++
	}

	category.Number = number

	return articles
}

// 获取 markdown 中的文章信息
func getArticleContent(path string) Article {
	article := Article{}

	// 读取文件
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		mylog.LogErr.Printf("read %q failure, %v", path, err)
		return article
	}

	// 字符串格式的文档
	str := string(bytes)

	// 将字符串按照 +++ 截取
	arr := strings.SplitN(str, "+++", 3)

	// 正常截取完是 3 部分，如果不是，就直接返回 nil
	if len(arr) != 3 {
		mylog.LogErr.Printf("file %q is incomplete, intercepted failure", path)
		return article
	}

	head := strings.TrimSpace(arr[1])
	body := strings.TrimSpace(arr[2])

	if head == "" || body == "" {
		mylog.LogErr.Printf("file %q head or body is empty", path)
		return article
	}

	// 解析头部的 toml
	if _, err := toml.Decode(head, &article); err != nil {
		mylog.LogErr.Printf("file %q parse toml head failure, %v", path, err)
		return article
	}

	article.Body = body

	return article
}

//下面三个方法是实现了 sort 接口，实现 Articles 的排序

// Len
func (a Articles) Len() int {
	return len(a)
}

// Swap 实现的 sort 接口
func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 进行排序比较， 第一排序是 UpdatedAt ，
// 第二排序是 CreatedAt， 第三排序是 Title
func (a Articles) Less(i, j int) bool {
	if a[i].UpdatedAt.After(a[j].UpdatedAt) {
		return true
	}

	if a[i].CreatedAt.After(a[j].CreatedAt) {
		return true
	}

	if a[i].Title > a[j].Title {
		return true
	}

	return false
}
