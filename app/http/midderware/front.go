package midderware

import (
	"path"
	"sort"
	"strings"

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
	pathValue := c.Param("name")

	for _, category := range categories {
		if category.OutLink || category.Number > 0 {
			// 如果当前参数路径和分类路径相同，就是激活状态
			if strings.ToLower(pathValue) == strings.ToLower(category.Path) ||
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
	name := c.Param("name")

	// 按照标签数量排序
	sort.Sort(&tags)

	for i := 0; i < len(tags); i++ {
		if strings.ToLower(name) == tags[i].Title {
			tags[i].Active = true
		} else {
			tags[i].Active = false
		}
	}

	c.Set("tags", tags)
	c.Next()
}
