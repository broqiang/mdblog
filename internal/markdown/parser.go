package markdown

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mdblog/internal/config"
	"mdblog/internal/data"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Parser Markdown解析器
type Parser struct {
	md goldmark.Markdown
}

// NewParser 创建新的Markdown解析器
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Linkify,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // 允许原始HTML
		),
	)

	return &Parser{
		md: md,
	}
}

// ParseFile 解析Markdown文件
func (m *Parser) ParseFile(filePath, postsDir string) (*data.Post, error) {
	// 读取文件内容
	content, err := readFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析Front Matter
	frontMatter, markdownContent, err := ParseFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("解析Front Matter失败: %w", err)
	}

	// 生成文章ID
	id := generateID(filePath, postsDir)

	// 获取分类
	category := getCategoryFromPath(filePath, postsDir)

	// 解析Markdown为HTML
	var buf bytes.Buffer
	if err := m.md.Convert([]byte(markdownContent), &buf); err != nil {
		return nil, fmt.Errorf("转换Markdown失败: %w", err)
	}

	// 生成摘要
	summary := GenerateSummary(markdownContent, config.SummaryLines)

	// 处理时间
	createTime := frontMatter.CreatedAt.Time
	if createTime.IsZero() {
		createTime = time.Now()
	}

	updateTime := frontMatter.UpdatedAt.Time
	if updateTime.IsZero() {
		updateTime = createTime
	}

	// 创建文章对象
	post := &data.Post{
		ID:         id,
		Title:      frontMatter.Title,
		Author:     frontMatter.Author,
		GitHubURL:  frontMatter.GitHubURL,
		Content:    markdownContent,
		HTML:       buf.String(),
		Summary:    summary,
		Category:   category,
		CreateTime: createTime,
		UpdateTime: updateTime,
		Tags:       frontMatter.Tags,
		FilePath:   filePath,
	}

	// 如果没有标题，从文件名生成
	if post.Title == "" {
		post.Title = generateTitleFromPath(filePath)
	}

	return post, nil
}

// ParseContent 解析Markdown内容
func (m *Parser) ParseContent(content string) (string, error) {
	var buf bytes.Buffer
	if err := m.md.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("转换Markdown失败: %w", err)
	}
	return buf.String(), nil
}

// readFile 读取文件内容
func readFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件 %s 失败: %w", filePath, err)
	}
	return string(data), nil
}

// generateID 根据文件路径生成文章ID
func generateID(filePath, postsDir string) string {
	// 移除扩展名和posts目录前缀
	relPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	if strings.HasPrefix(relPath, postsDir+"/") {
		relPath = strings.TrimPrefix(relPath, postsDir+"/")
	}
	return relPath
}

// getCategoryFromPath 从文件路径获取分类
func getCategoryFromPath(filePath, postsDir string) string {
	relPath := filePath
	if strings.HasPrefix(relPath, postsDir+"/") {
		relPath = strings.TrimPrefix(relPath, postsDir+"/")
	}

	parts := strings.Split(relPath, "/")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

// generateTitleFromPath 从文件路径生成标题
func generateTitleFromPath(filePath string) string {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// 将下划线和连字符替换为空格
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	// 首字母大写
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}
