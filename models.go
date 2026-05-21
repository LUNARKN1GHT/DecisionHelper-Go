package main

// Criterion 决策标准，包含权重（1-5）
type Criterion struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

// Score 某个选项在某个标准上的评分（1-5）
type Score struct {
	Option      string `json:"option"`
	CriterionID string `json:"criterion_id"`
	Value       int    `json:"value"`
}

// Decision 一条决策记录
type Decision struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	CreatedAt string      `json:"created_at"`
	Options   []string    `json:"options"`
	Criteria  []Criterion `json:"criteria"`
	Scores    []Score     `json:"scores"`
}

// OptionResult 计算结果：选项、加权得分、排名
type OptionResult struct {
	Option string  `json:"option"`
	Score  float64 `json:"score"`
	Rank   int     `json:"rank"`
}
