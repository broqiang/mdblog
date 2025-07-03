package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mdblog/internal/config"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// verifyGiteeSignature 验证 Gitee Webhook 签名
func verifyGiteeSignature(signature, timestamp, secret string) bool {
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	stringToSign := fmt.Sprintf("%d\n%s", timestampInt, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)
	encodedSign := base64.StdEncoding.EncodeToString(signData)
	urlEncodedSign := url.PathEscape(encodedSign)
	return urlEncodedSign == signature
}

// handleGiteeWebhook 处理 Gitee Webhook 请求
func (s *Server) handleGiteeWebhook(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "方法不被允许"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
		return
	}
	giteeToken := c.GetHeader("X-Gitee-Token")
	giteeTimestamp := c.GetHeader("X-Gitee-Timestamp")
	if giteeToken == "" {
		giteeToken = c.Query("sign")
	}
	if giteeTimestamp == "" {
		giteeTimestamp = c.Query("timestamp")
	}
	if giteeToken == "" || giteeTimestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少签名或时间戳"})
		return
	}
	timestampInt, err := strconv.ParseInt(giteeTimestamp, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的时间戳格式"})
		return
	}
	currentTime := time.Now().Unix() * 1000
	if diff := currentTime - timestampInt; diff > 3600000 || diff < -3600000 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "时间戳超出允许范围"})
		return
	}
	if !verifyGiteeSignature(giteeToken, giteeTimestamp, config.WebhookSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
		return
	}
	var payload GiteeWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 解析失败"})
		return
	}
	if payload.Ref != "refs/heads/main" {
		c.JSON(http.StatusOK, gin.H{"message": "忽略非 main 分支的推送"})
		return
	}
	postsDir := s.manager.GetPostsDir()
	cmd := exec.Command("git", "pull")
	cmd.Dir = postsDir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Git pull 失败", "details": string(output)})
		return
	}
	if reloadErr := s.reloadData(); reloadErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重新加载数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "同步成功",
		"branch":     "main",
		"commits":    len(payload.Commits),
		"repository": payload.Repository.Name,
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GiteeWebhookPayload Gitee Webhook 载荷结构体
// 只保留必要字段
type GiteeWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Commits []struct {
		ID string `json:"id"`
	} `json:"commits"`
}
