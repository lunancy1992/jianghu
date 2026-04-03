package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DeepSeekClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewDeepSeekClient(apiKey, baseURL, model string) *DeepSeekClient {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	if model == "" {
		model = "deepseek-chat"
	}
	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// OpenAI-compatible request/response types.

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// SummaryResult holds the AI-generated summary.
type SummaryResult struct {
	OneLiner    string   `json:"one_liner"`
	Category    string   `json:"category"`
	Keywords    []string `json:"keywords"`
	Sensitivity int      `json:"sensitivity"`
}

// ModerationResult holds the AI moderation outcome.
type ModerationResult struct {
	Passed     bool   `json:"passed"`
	Rejected   bool   `json:"rejected"`
	NeedReview bool   `json:"need_review"`
	Reason     string `json:"reason"`
}

// Summarize generates an objective summary of a news article.
func (c *DeepSeekClient) Summarize(ctx context.Context, title, content string) (*SummaryResult, error) {
	prompt := fmt.Sprintf(`你是一个客观中立的新闻摘要生成器。请分析以下新闻，返回JSON格式：
{
  "one_liner": "一句话概括（不超过50字）",
  "category": "分类（国际/国内/科技/财经/社会/文化/其他）",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "sensitivity": 0-10的敏感度评分
}

标题：%s
内容：%s

只返回JSON，不要其他文字。`, title, content)

	resp, err := c.chat(ctx, prompt, 0.3, 500)
	if err != nil {
		return nil, err
	}

	var result SummaryResult
	cleaned := extractJSON(resp)
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("parse summary response: %w (raw: %s)", err, resp)
	}
	return &result, nil
}

// Moderate checks if content is appropriate.
func (c *DeepSeekClient) Moderate(ctx context.Context, content string) (*ModerationResult, error) {
	prompt := fmt.Sprintf(`你是一个内容审核员。请审核以下用户评论内容，返回JSON格式：
{
  "passed": true/false,
  "rejected": true/false,
  "need_review": true/false,
  "reason": "原因说明"
}

规则：
- 不含人身攻击、仇恨言论、色情、违法内容则passed=true
- 明确违规则rejected=true
- 模糊地带则need_review=true
- 三个字段互斥，只能有一个为true

评论内容：%s

只返回JSON，不要其他文字。`, content)

	resp, err := c.chat(ctx, prompt, 0.1, 300)
	if err != nil {
		return nil, err
	}

	var result ModerationResult
	cleaned := extractJSON(resp)
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("parse moderation response: %w (raw: %s)", err, resp)
	}
	return &result, nil
}

// extractJSON strips markdown code fences and whitespace from AI responses.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	// Remove ```json ... ``` or ``` ... ```
	if strings.HasPrefix(s, "```") {
		// Find end of first line (```json or ```)
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		// Remove trailing ```
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
	}
	return strings.TrimSpace(s)
}

func (c *DeepSeekClient) chat(ctx context.Context, prompt string, temperature float64, maxTokens int) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("deepseek api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deepseek api error %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
