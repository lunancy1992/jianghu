package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClaudeClient implements the same Summarize/Moderate interface using Anthropic Messages API.
type ClaudeClient struct {
	apiKey string
	model  string
	client *http.Client
}

func NewClaudeClient(apiKey, model string) *ClaudeClient {
	if model == "" {
		model = "claude-haiku-4-5-20251001"
	}
	return &ClaudeClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Anthropic Messages API types.

type claudeRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	Messages  []claudeMessage  `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (c *ClaudeClient) Summarize(ctx context.Context, title, content string) (*SummaryResult, error) {
	prompt := fmt.Sprintf(`你是一个客观中立的新闻摘要生成器。请分析以下新闻，返回JSON格式：
{
  "one_liner": "一句话概括（不超过50字）",
  "category": "分类（政治/经济/科技/社会/国际/娱乐/体育/其他）",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "sensitivity": 0-10的敏感度评分
}

标题：%s
内容：%s

只返回JSON，不要其他文字。`, title, truncate(content, 3000))

	resp, err := c.chat(ctx, prompt, 500)
	if err != nil {
		return nil, err
	}

	var result SummaryResult
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, fmt.Errorf("parse summary response: %w (raw: %s)", err, resp)
	}
	return &result, nil
}

func (c *ClaudeClient) Moderate(ctx context.Context, content string) (*ModerationResult, error) {
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

	resp, err := c.chat(ctx, prompt, 300)
	if err != nil {
		return nil, err
	}

	var result ModerationResult
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, fmt.Errorf("parse moderation response: %w (raw: %s)", err, resp)
	}
	return &result, nil
}

func (c *ClaudeClient) chat(ctx context.Context, prompt string, maxTokens int) (string, error) {
	reqBody := claudeRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("claude api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude api error %d: %s", resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return claudeResp.Content[0].Text, nil
}

// truncate limits text to maxRunes to avoid exceeding token limits.
func truncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes])
}
