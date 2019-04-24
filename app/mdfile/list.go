// Package mdfile 是用来解析处理 markdown 文件的包
package mdfile

import (
	"errors"
	"path/filepath"

	"github.com/broqiang/mdblog/app/helper"

	"github.com/BurntSushi/toml"
	"github.com/broqiang/mdblog/app/config"
)

const mdDir = "resource"

// Model 博客内容的实例
var Model = new()

// List 博客列表
type List interface {
	CategoriesAll() Categories
	TagsAll() Tags
	ArticlesAll() Articles
}

// Categories 是文章分类的切片（数组）
type Categories []Category

// Category 文章分类
type Category struct {
	// 分类的名称，用来做头部菜单的文本
	Title string

	// 分类下文章的数量
	Number int

	// 分类的文件目录， 当 OutLink 为 true 时，这里应该是一个外链，
	// 否则是 markdown 文档保存的子目录
	Path string

	// 分类的描述
	Description string

	// 当是一个外联时，就不会再去获取对应的 markdown 文档
	OutLink bool `toml:"out_link"`

	// 是否是激活状态
	Active bool
}

// Initialize 初始化数据
func new() List {
	// 暂时是保存到内存，后面计划支持其他存储方式
	l := newListMap()

	return l
}

// 获取分类配置文件的位置
func getCategoriesPath() string {
	return filepath.Join(config.Root, "config", "categories.toml")
}

// 获取 markdown 文档的根目录
func getRootPath() string {
	return filepath.Join(config.Root, config.Cfg.MarkdownDir)
}

// 解析分类配置文件
func parseCategories() Categories {
	temp := struct {
		Category []Category
	}{}

	if _, err := toml.DecodeFile(getCategoriesPath(), &temp); err != nil {
		helper.Panicf("cannot parse categories.toml config file. %v", err)
	}

	if temp.Category == nil {
		helper.PanicErr(errors.New("have not been defined category"))
	}

	return temp.Category
}
