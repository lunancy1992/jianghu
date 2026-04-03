package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/config"
)

type Provider interface {
	Send(ctx context.Context, phone, code string) error
}

// AliyunSMS 阿里云短信服务
type AliyunSMS struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

func NewAliyunSMS(cfg *config.Config) *AliyunSMS {
	return &AliyunSMS{
		AccessKeyID:     cfg.Auth.SMSAccessKey,
		AccessKeySecret: cfg.Auth.SMSAccessSecret,
		SignName:        cfg.Auth.SMSSignName,
		TemplateCode:    cfg.Auth.SMSTemplateCode,
	}
}

func (s *AliyunSMS) Send(ctx context.Context, phone, code string) error {
	params := map[string]string{
		"AccessKeyId":      s.AccessKeyID,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     phone,
		"RegionId":         "cn-hangzhou",
		"SignName":         s.SignName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"SignatureVersion": "1.0",
		"TemplateCode":     s.TemplateCode,
		"TemplateParam":    fmt.Sprintf(`{"code":"%s"}`, code),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2017-05-25",
	}

	// 构造签名字符串
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts,
			fmt.Sprintf("%s=%s", specialURLEncode(k), specialURLEncode(params[k])))
	}
	canonicalizedQueryString := strings.Join(queryParts, "&")

	stringToSign := "GET&" + specialURLEncode("/") + "&" + specialURLEncode(canonicalizedQueryString)

	// 计算签名
	mac := hmac.New(sha1.New, []byte(s.AccessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	params["Signature"] = signature

	// 发送请求
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}

	reqURL := "https://dysmsapi.aliyuncs.com/?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Code != "OK" {
		return fmt.Errorf("sms send failed: %s", result.Message)
	}

	return nil
}

// specialURLEncode 阿里云特殊的 URL 编码
func specialURLEncode(s string) string {
	encoded := url.QueryEscape(s)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

// StubSMS 开发环境用的 Stub 实现
type StubSMS struct{}

func NewStubSMS() *StubSMS {
	return &StubSMS{}
}

func (s *StubSMS) Send(ctx context.Context, phone, code string) error {
	fmt.Printf("[SMS STUB] Sending code %s to phone %s\n", code, phone)
	return nil
}
