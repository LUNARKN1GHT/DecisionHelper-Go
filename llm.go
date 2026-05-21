package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ── OpenAI 兼容请求/响应结构 ───────────────────────────────────────────────

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ── AI 建议的数据结构（暴露给前端）────────────────────────────────────────

type SuggestedCriterion struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Reason string `json:"reason"`
}

type SuggestedScore struct {
	Option    string `json:"option"`
	Criterion string `json:"criterion"`
	Value     int    `json:"value"`
	Reason    string `json:"reason"`
}

// ── 核心调用 ───────────────────────────────────────────────────────────────

func callLLM(cfg LLMConfig, system, user string) (string, error) {
	req := chatRequest{
		Model: cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Temperature: 0.3,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := strings.TrimRight(cfg.BaseURL, "/") + "/v1/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("请求失败，请检查网络或 API 地址：%w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var cr chatResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return "", fmt.Errorf("解析响应失败：%w", err)
	}
	if cr.Error != nil {
		return "", fmt.Errorf("API 错误：%s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("API 返回了空结果")
	}
	return cr.Choices[0].Message.Content, nil
}

// extractJSON 去除 LLM 输出中常见的 markdown 代码块包裹
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "```json"); idx != -1 {
		s = s[idx+7:]
	} else if idx := strings.Index(s, "```"); idx != -1 {
		s = s[idx+3:]
	}
	if idx := strings.LastIndex(s, "```"); idx != -1 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ── AI 建议标准 ────────────────────────────────────────────────────────────

func suggestCriteriaLLM(cfg LLMConfig, d *Decision) ([]SuggestedCriterion, error) {
	system := `你是一个专业的决策分析顾问。用户会告诉你一个决策场景，你需要推荐5个最关键的评价标准。
严格以JSON数组格式返回，不要有任何其他文字，格式示例：
[{"name":"标准名称","weight":4,"reason":"简短理由（不超过20字）"}]
weight 取值 1-5，代表重要程度。`

	options := strings.Join(d.Options, "、")
	if options == "" {
		options = "（暂无）"
	}
	user := fmt.Sprintf("决策标题：%s\n候选选项：%s\n请推荐5个评价标准。", d.Title, options)

	content, err := callLLM(cfg, system, user)
	if err != nil {
		return nil, err
	}
	var result []SuggestedCriterion
	if err := json.Unmarshal([]byte(extractJSON(content)), &result); err != nil {
		return nil, fmt.Errorf("解析标准建议失败，请重试")
	}
	for i := range result {
		result[i].Weight = clampInt(result[i].Weight, 1, 5)
	}
	return result, nil
}

// ── AI 辅助评分 ────────────────────────────────────────────────────────────

func suggestScoresLLM(cfg LLMConfig, d *Decision) ([]SuggestedScore, error) {
	system := `你是一个专业的决策分析顾问。根据用户提供的决策信息，对每个"选项×标准"组合给出建议评分。
严格以JSON数组格式返回，不要有任何其他文字，格式示例：
[{"option":"选项名","criterion":"标准名","value":4,"reason":"简短理由（不超过20字）"}]
value 取值 1-5，1=很差，5=很好。`

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("决策标题：%s\n选项：%s\n标准：\n",
		d.Title, strings.Join(d.Options, "、")))
	for _, c := range d.Criteria {
		sb.WriteString(fmt.Sprintf("- %s（重要性%d）\n", c.Name, c.Weight))
	}
	sb.WriteString("请对所有选项和标准的组合进行评分。")

	content, err := callLLM(cfg, system, sb.String())
	if err != nil {
		return nil, err
	}
	var result []SuggestedScore
	if err := json.Unmarshal([]byte(extractJSON(content)), &result); err != nil {
		return nil, fmt.Errorf("解析评分建议失败，请重试")
	}
	for i := range result {
		result[i].Value = clampInt(result[i].Value, 1, 5)
	}
	return result, nil
}

// ── AI 分析结果 ────────────────────────────────────────────────────────────

func analyzeResultsLLM(cfg LLMConfig, d *Decision, results []OptionResult) (string, error) {
	system := `你是一个专业的决策分析顾问。根据用户的决策矩阵分析结果，用中文给出一段简洁分析（150字以内）：
指出推荐选项及核心理由，并提示需权衡的关键点或潜在风险。直接输出正文，不要使用标题、列表或 Markdown。`

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("决策：%s\n\n评价标准：", d.Title))
	for _, c := range d.Criteria {
		sb.WriteString(fmt.Sprintf("%s（权重%d）、", c.Name, c.Weight))
	}
	sb.WriteString("\n\n加权得分排名：\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("第%d名：%s（%.2f 分）\n", r.Rank, r.Option, r.Score))
	}

	return callLLM(cfg, system, sb.String())
}
