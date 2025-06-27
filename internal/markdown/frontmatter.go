package markdown

import (
	"strings"
	"time"

	"mdblog/internal/data"

	"gopkg.in/yaml.v3"
)

// ParseFrontMatter 解析Front Matter
func ParseFrontMatter(content string) (*data.FrontMatter, string, error) {
	lines := strings.Split(content, "\n")

	// 检查是否有Front Matter
	if len(lines) < 3 || !strings.HasPrefix(lines[0], "---") {
		return &data.FrontMatter{}, content, nil
	}

	// 查找结束标记
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return &data.FrontMatter{}, content, nil
	}

	// 提取Front Matter内容
	frontMatterLines := lines[1:endIndex]
	frontMatterContent := strings.Join(frontMatterLines, "\n")

	// 解析YAML
	var fm data.FrontMatter
	if err := yaml.Unmarshal([]byte(frontMatterContent), &fm); err != nil {
		return nil, "", err
	}

	// 提取Markdown内容
	markdownContent := strings.Join(lines[endIndex+1:], "\n")

	return &fm, markdownContent, nil
}

// GenerateSummary 生成文章摘要
func GenerateSummary(content string, maxLines int) string {
	lines := strings.Split(content, "\n")

	// 过滤空行和代码块
	var summaryLines []string
	inCodeBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行
		if line == "" {
			continue
		}

		// 检查代码块标记
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// 跳过代码块内的内容
		if inCodeBlock {
			continue
		}

		// 跳过标题行
		if strings.HasPrefix(line, "#") {
			continue
		}

		summaryLines = append(summaryLines, line)

		if len(summaryLines) >= maxLines {
			break
		}
	}

	summary := strings.Join(summaryLines, "\n")

	// 如果内容太长，截断并添加省略号
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	return summary
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	// 尝试多种时间格式
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	// 如果都解析失败，返回当前时间
	return time.Now(), nil
}
