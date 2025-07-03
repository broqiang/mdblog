package server

import (
	"html/template"
	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// initTemplates 初始化模板
func (s *Server) initTemplates() {
	// 定义模板函数
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"ceil": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return (a + b - 1) / b // 向上取整
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case map[string]int:
				return len(val)
			case []string:
				return len(val)
			default:
				return 0
			}
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// 优先使用嵌入的文件系统
	if s.embeddedAssets != nil {
		log.Println("使用嵌入的模板文件系统")

		// 从嵌入的文件系统创建模板子文件系统
		templatesSubFS, err := fs.Sub(*s.embeddedAssets, "web/templates")
		if err != nil {
			log.Fatalf("无法创建嵌入模板子文件系统: %v", err)
		}

		// 解析所有模板文件，需要指定具体的路径
		tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesSubFS,
			"layouts/head.html",
			"layouts/header.html",
			"layouts/footer.html",
			"layouts/js.html",
			"posts/index.html",
			"posts/detail.html",
			"posts/category.html",
			"404.html",
			"search.html")
		if err != nil {
			log.Fatalf("解析嵌入模板失败: %v", err)
		}

		// 设置模板到 Gin 引擎
		s.engine.SetHTMLTemplate(tmpl)
	} else {
		log.Println("使用文件系统模板")

		// 检查模板目录是否存在
		templatesDir := "web/templates"
		if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
			log.Fatalf("模板目录不存在: %s。请确保在项目根目录运行或使用嵌入模式", templatesDir)
		}

		// 使用 Gin 的模板引擎，加载所有子目录的模板文件
		s.engine.SetFuncMap(funcMap)
		s.engine.LoadHTMLGlob("web/templates/*/*")
	}
}

// renderTemplate 渲染模板
func (s *Server) renderTemplate(c *gin.Context, templateName string, data interface{}) {
	c.HTML(200, templateName, data)
}
