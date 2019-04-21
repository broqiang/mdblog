// Package markdown 是用于处理 markdown 文档转换的
package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"gopkg.in/russross/blackfriday.v1"
)

// HTML 将 markdown 文档转换成 html 文件
func HTML(input string) template.HTML {

	// 因为默认字符串会被模板转译，所以返回一个 template.HTML
	// 就可以叫 HTML 原样输出了

	htmlContent := markdown([]byte(input))

	return template.HTML(htmlContent)
}

func toHTML(text []byte) (htmlContent []byte) {
	unsafe := markdown(text)
	htmlContent = policy.SanitizeBytes(unsafe)

	return
}

// 自定义的 blackfriday ，为了支持生成锚点
// 下面的主要内容是从 https://github.com/shurcooL/github_flavored_markdown 扒出来的
func markdown(text []byte) []byte {
	renderer := &renderer{Html: blackfriday.HtmlRenderer(markdownHTMLFlags, "", "").(*blackfriday.Html)}

	return blackfriday.MarkdownOptions(text, renderer, blackfriday.Options{
		Extensions: markdownExtensions})
}

// policy for GitHub Flavored Markdown-like sanitization.
var policy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("checked", "disabled").Matching(regexp.MustCompile(`^$`)).OnElements("input")
	p.AllowDataURIImages()
	return p
}()

// 定义 HTML 渲染器的配置选项
const markdownHTMLFlags = 0 |
	blackfriday.HTML_USE_XHTML |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_SMARTYPANTS_DASHES |
	blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
	blackfriday.HTML_FOOTNOTE_RETURN_LINKS

// 定义 markdown 扩展，其实就是复制的 commonExtensions
const markdownExtensions = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_HEADER_IDS |
	blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
	blackfriday.EXTENSION_DEFINITION_LISTS

type renderer struct {
	*blackfriday.Html
}

// GitHub Flavored Markdown heading with clickable and hidden anchor.
func (*renderer) Header(out *bytes.Buffer, text func() bool, level int, _ string) {
	marker := out.Len()
	doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}

	textHTML := out.String()[marker:]
	out.Truncate(marker)

	// Extract text content of the heading.
	var textContent string
	if node, err := html.Parse(strings.NewReader(textHTML)); err == nil {
		textContent = extractText(node)
	} else {
		// Failed to parse HTML (probably can never happen), so just use the whole thing.
		textContent = html.UnescapeString(textHTML)
	}
	anchorName := create(textContent)

	out.WriteString(fmt.Sprintf(`<h%d><a name="%s" class="anchor" href="#%s" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>`, level, anchorName, anchorName))
	out.WriteString(textHTML)
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func extractText(n *html.Node) string {
	var out string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			out += c.Data
		} else {
			out += extractText(c)
		}
	}
	return out
}

// Create returns a sanitized anchor name for the given text.
func create(text string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}
