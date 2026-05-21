package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDecisions() ([]Decision, error) {
	return loadDecisions()
}

func (a *App) GetDecision(id string) (*Decision, error) {
	decisions, err := loadDecisions()
	if err != nil {
		return nil, err
	}
	for i := range decisions {
		if decisions[i].ID == id {
			return &decisions[i], nil
		}
	}
	return nil, fmt.Errorf("decision not found: %s", id)
}

func (a *App) CreateDecision(title string) (*Decision, error) {
	decisions, err := loadDecisions()
	if err != nil {
		return nil, err
	}
	d := Decision{
		ID:        uuid.New().String(),
		Title:     title,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05"),
		Options:   []string{},
		Criteria:  []Criterion{},
		Scores:    []Score{},
	}
	decisions = append(decisions, d)
	if err := saveDecisions(decisions); err != nil {
		return nil, err
	}
	return &d, nil
}

func (a *App) DeleteDecision(id string) error {
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	filtered := make([]Decision, 0, len(decisions))
	for _, d := range decisions {
		if d.ID != id {
			filtered = append(filtered, d)
		}
	}
	return saveDecisions(filtered)
}

func (a *App) UpdateDecisionTitle(id, title string) error {
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == id {
			decisions[i].Title = title
			return saveDecisions(decisions)
		}
	}
	return fmt.Errorf("decision not found: %s", id)
}

func (a *App) AddOption(decisionID, option string) error {
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			for _, o := range decisions[i].Options {
				if o == option {
					return fmt.Errorf("option already exists: %s", option)
				}
			}
			decisions[i].Options = append(decisions[i].Options, option)
			return saveDecisions(decisions)
		}
	}
	return fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) RemoveOption(decisionID, option string) error {
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			filtered := make([]string, 0)
			for _, o := range decisions[i].Options {
				if o != option {
					filtered = append(filtered, o)
				}
			}
			decisions[i].Options = filtered
			scores := make([]Score, 0)
			for _, s := range decisions[i].Scores {
				if s.Option != option {
					scores = append(scores, s)
				}
			}
			decisions[i].Scores = scores
			return saveDecisions(decisions)
		}
	}
	return fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) AddCriterion(decisionID, name string, weight int) (*Criterion, error) {
	if weight < 1 || weight > 5 {
		return nil, fmt.Errorf("weight must be between 1 and 5")
	}
	decisions, err := loadDecisions()
	if err != nil {
		return nil, err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			c := Criterion{
				ID:     uuid.New().String(),
				Name:   name,
				Weight: weight,
			}
			decisions[i].Criteria = append(decisions[i].Criteria, c)
			if err := saveDecisions(decisions); err != nil {
				return nil, err
			}
			return &c, nil
		}
	}
	return nil, fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) UpdateCriterion(decisionID, criterionID, name string, weight int) error {
	if weight < 1 || weight > 5 {
		return fmt.Errorf("weight must be between 1 and 5")
	}
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			for j := range decisions[i].Criteria {
				if decisions[i].Criteria[j].ID == criterionID {
					decisions[i].Criteria[j].Name = name
					decisions[i].Criteria[j].Weight = weight
					return saveDecisions(decisions)
				}
			}
			return fmt.Errorf("criterion not found: %s", criterionID)
		}
	}
	return fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) RemoveCriterion(decisionID, criterionID string) error {
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			filtered := make([]Criterion, 0)
			for _, c := range decisions[i].Criteria {
				if c.ID != criterionID {
					filtered = append(filtered, c)
				}
			}
			decisions[i].Criteria = filtered
			scores := make([]Score, 0)
			for _, s := range decisions[i].Scores {
				if s.CriterionID != criterionID {
					scores = append(scores, s)
				}
			}
			decisions[i].Scores = scores
			return saveDecisions(decisions)
		}
	}
	return fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) SetScore(decisionID, option, criterionID string, value int) error {
	if value < 1 || value > 5 {
		return fmt.Errorf("score must be between 1 and 5")
	}
	decisions, err := loadDecisions()
	if err != nil {
		return err
	}
	for i := range decisions {
		if decisions[i].ID == decisionID {
			for j := range decisions[i].Scores {
				if decisions[i].Scores[j].Option == option && decisions[i].Scores[j].CriterionID == criterionID {
					decisions[i].Scores[j].Value = value
					return saveDecisions(decisions)
				}
			}
			decisions[i].Scores = append(decisions[i].Scores, Score{
				Option:      option,
				CriterionID: criterionID,
				Value:       value,
			})
			return saveDecisions(decisions)
		}
	}
	return fmt.Errorf("decision not found: %s", decisionID)
}

func (a *App) GetResults(decisionID string) ([]OptionResult, error) {
	d, err := a.GetDecision(decisionID)
	if err != nil {
		return nil, err
	}

	totalWeight := 0
	for _, c := range d.Criteria {
		totalWeight += c.Weight
	}

	results := make([]OptionResult, 0, len(d.Options))
	for _, option := range d.Options {
		if totalWeight == 0 {
			results = append(results, OptionResult{Option: option, Score: 0})
			continue
		}
		weightedSum := 0.0
		for _, c := range d.Criteria {
			for _, s := range d.Scores {
				if s.Option == option && s.CriterionID == c.ID {
					weightedSum += float64(s.Value) * float64(c.Weight)
					break
				}
			}
		}
		results = append(results, OptionResult{
			Option: option,
			Score:  weightedSum / float64(totalWeight),
		})
	}

	// 按加权得分降序排序，分配名次
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	for i := range results {
		results[i].Rank = i + 1
	}

	return results, nil
}

// ── LLM 配置 ──────────────────────────────────────────────────────────────

func (a *App) GetLLMConfig() (LLMConfig, error) {
	return loadLLMConfig()
}

func (a *App) SaveLLMConfig(cfg LLMConfig) error {
	return saveLLMConfig(cfg)
}

// ── AI 辅助功能 ────────────────────────────────────────────────────────────

func (a *App) SuggestCriteria(decisionID string) ([]SuggestedCriterion, error) {
	cfg, err := loadLLMConfig()
	if err != nil {
		return nil, err
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("请先在设置中配置 API Key")
	}
	d, err := a.GetDecision(decisionID)
	if err != nil {
		return nil, err
	}
	return suggestCriteriaLLM(cfg, d)
}

func (a *App) SuggestScores(decisionID string) ([]SuggestedScore, error) {
	cfg, err := loadLLMConfig()
	if err != nil {
		return nil, err
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("请先在设置中配置 API Key")
	}
	d, err := a.GetDecision(decisionID)
	if err != nil {
		return nil, err
	}
	if len(d.Options) == 0 || len(d.Criteria) == 0 {
		return nil, fmt.Errorf("请先添加选项和标准")
	}
	return suggestScoresLLM(cfg, d)
}

func (a *App) AnalyzeResults(decisionID string) (string, error) {
	cfg, err := loadLLMConfig()
	if err != nil {
		return "", err
	}
	if cfg.APIKey == "" {
		return "", fmt.Errorf("请先在设置中配置 API Key")
	}
	d, err := a.GetDecision(decisionID)
	if err != nil {
		return "", err
	}
	results, err := a.GetResults(decisionID)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("没有可分析的结果，请先完成评分")
	}
	return analyzeResultsLLM(cfg, d, results)
}
