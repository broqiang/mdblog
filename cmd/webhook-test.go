package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// WebhookTestPayload 测试用的 Webhook 载荷
type WebhookTestPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name       string `json:"name"`
		FullName   string `json:"full_name"`
		CloneURL   string `json:"clone_url"`
		SSHURL     string `json:"ssh_url"`
		GitHTTPURL string `json:"git_http_url"`
		GitSSHURL  string `json:"git_ssh_url"`
	} `json:"repository"`
	Commits []struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		ID        string    `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"head_commit"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
}

// generateGiteeSignature 按照 Gitee 官方文档生成签名
func generateGiteeSignature(timestamp, secret string) string {
	// 将时间戳字符串转换为 int64
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		fmt.Printf("签名生成错误: 无法解析时间戳 %s: %v", timestamp, err)
		return ""
	}

	// Step1: 构造签名字符串: timestamp + "\n" + secret (使用 int64 格式)
	stringToSign := fmt.Sprintf("%d\n%s", timestampInt, secret)

	// Step2: 使用 HmacSHA256 算法计算签名
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)

	// Step3: Base64 编码
	encodedSign := base64.StdEncoding.EncodeToString(signData)

	// Step4: URL 编码 (使用 PathEscape 而不是 QueryEscape)
	urlEncodedSign := url.PathEscape(encodedSign)

	return urlEncodedSign
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: go run cmd/webhook-test.go <webhook_url> <secret_key>")
		fmt.Println("示例: go run cmd/webhook-test.go http://localhost:8091/webhook/gitee c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5")
		os.Exit(1)
	}

	webhookURL := os.Args[1]
	secretKey := os.Args[2]

	// 创建测试载荷
	payload := WebhookTestPayload{
		Ref: "refs/heads/main",
		Repository: struct {
			Name       string `json:"name"`
			FullName   string `json:"full_name"`
			CloneURL   string `json:"clone_url"`
			SSHURL     string `json:"ssh_url"`
			GitHTTPURL string `json:"git_http_url"`
			GitSSHURL  string `json:"git_ssh_url"`
		}{
			Name:       "mdblog-posts",
			FullName:   "broqiang/mdblog-posts",
			CloneURL:   "https://gitee.com/broqiang/mdblog-posts.git",
			SSHURL:     "git@gitee.com:broqiang/mdblog-posts.git",
			GitHTTPURL: "https://gitee.com/broqiang/mdblog-posts.git",
			GitSSHURL:  "git@gitee.com:broqiang/mdblog-posts.git",
		},
		Commits: []struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
			Added    []string `json:"added"`
			Removed  []string `json:"removed"`
			Modified []string `json:"modified"`
		}{
			{
				ID:        "abc123def456",
				Message:   "更新博客文章",
				Timestamp: time.Now(),
				Author: struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					Name:  "BroQiang",
					Email: "broqiang@example.com",
				},
				Added:    []string{"posts/test/new-article.md"},
				Removed:  []string{},
				Modified: []string{"posts/about.md"},
			},
		},
		HeadCommit: struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		}{
			ID:        "abc123def456",
			Message:   "更新博客文章",
			Timestamp: time.Now(),
			Author: struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			}{
				Name:  "BroQiang",
				Email: "broqiang@example.com",
			},
		},
		Pusher: struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Name:  "BroQiang",
			Email: "broqiang@example.com",
		},
	}

	// 序列化载荷
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("序列化载荷失败: %v\n", err)
		os.Exit(1)
	}

	// 生成时间戳（毫秒）
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	// 按照 Gitee 官方文档生成签名
	signature := generateGiteeSignature(timestamp, secretKey)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		os.Exit(1)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gitee-Token", signature)
	req.Header.Set("X-Gitee-Timestamp", timestamp)
	req.Header.Set("X-Gitee-Event", "Push Hook")
	req.Header.Set("User-Agent", "Git-Webhook/1.0")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	fmt.Printf("发送测试 Webhook 到: %s\n", webhookURL)
	fmt.Printf("时间戳: %s\n", timestamp)
	fmt.Printf("签名: %s\n", signature)
	fmt.Printf("签名详情:\n")
	fmt.Printf("  签名字符串: %s\\n%s\n", timestamp, secretKey)

	// 显示签名计算过程
	stringToSign := timestamp + "\n" + secretKey
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)
	encodedSign := base64.StdEncoding.EncodeToString(signData)
	fmt.Printf("  Base64编码: %s\n", encodedSign)
	fmt.Printf("  URL编码: %s\n", signature)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		os.Exit(1)
	}

	// 打印结果
	fmt.Printf("响应状态: %s\n", resp.Status)
	fmt.Printf("响应内容: %s\n", string(body))

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Webhook 测试成功！")
	} else {
		fmt.Println("❌ Webhook 测试失败！")
		os.Exit(1)
	}
}
