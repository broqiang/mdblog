package midderware

import (
	"path"
	"sort"

	"github.com/broqiang/mdblog/app/config"

	"github.com/broqiang/mdblog/app/mdfile"
	"github.com/gin-gonic/gin"
)

// Sites 站点的一些公共信息配置
func Sites(c *gin.Context) {
	c.Set("cfg", config.Cfg)

	c.Next()
}

// Navigation 导航栏用到的中间件
func Navigation(c *gin.Context) {
	categories := mdfile.Model.CategoriesAll()
	cates := make(mdfile.Categories, 0)
	pathValue := c.Param("cname")

	for _, category := range categories {
		if category.OutLink || category.Number > 0 {
			// 如果当前参数路径和分类路径相同，就是激活状态
			if pathValue == category.Path ||
				c.Request.URL.Path == category.Path {
				category.Active = true
			}

			if !category.OutLink {
				category.Path = path.Join("/c", category.Path)
			}

			cates = append(cates, category)
		}
	}

	// 设置分类数据
	c.Set("categories", cates)
	c.Next()
}

// Tags 是右侧标签中用到的数据
func Tags(c *gin.Context) {

	tags := mdfile.Model.TagsAll()

	// 按照标签数量排序
	sort.Sort(&tags)

	c.Set("tags", tags)
	c.Next()
}
