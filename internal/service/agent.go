package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yumosx/agent/internal/domain"
	"github.com/yumosx/agent/internal/domain/params"
	"github.com/yumosx/agent/internal/service/llm"
	"regexp"
	"strings"
)

const (
	NO_STARTED  = "no_started"
	IN_PROGRESS = "in_progress"
	COMPLETED   = "completed"
	BLOCKED     = "blocked"
)

type AgentService struct {
	Id       string
	handler  *llm.Handler
	plan     *domain.Plan
	executor *PlanExecutor
}

func NewPlanService(handler *llm.Handler, executor *PlanExecutor) *AgentService {
	return &AgentService{handler: handler, executor: executor, plan: &domain.Plan{Id: "1"}}
}

func (p *AgentService) Plan(ctx context.Context, s string) (string, error) {
	var req domain.LLMRequest

	req.SystemContent = `You are a planning assistant. Create a concise, actionable plan with clear steps. 
Focus on key milestones rather than detailed sub-steps. 
Optimize for clarity and efficiency.`

	req.Msgs = []domain.Msg{
		{Role: domain.USER, Content: fmt.Sprintf("Create a reasonable plan with clear steps to accomplish the task: %s", s)}}

	req.Tools = []domain.Tool{p.newPlanTool()}

	err := p.createInitPlan(ctx, req)
	if err != nil {
		return "", err
	}
	plan := p.formatPlan()
	return plan, nil
}

func (p *AgentService) createInitPlan(ctx context.Context, req domain.LLMRequest) error {
	resp, err := p.handler.Invoke(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.ToolCalls) == 0 {
		return errors.New("LLM 返回的 toolCalls 为空")
	}

	for _, t := range resp.ToolCalls {
		if t.Function.Name == "planning" {
			return p.initPlanWithArgs(t.Function.Arguments)
		}
	}

	return errors.New("LLM 返回的 Function name 非法")
}

func (p *AgentService) Execute(ctx context.Context) error {
	var err error
	for {
		var (
			index int
			step  string
		)
		index, step, err = p.getStepInfo()
		if err != nil {
			return err
		}
		err = p.executeStep(ctx, p.executor, index, step)
		if err != nil {
			return err
		}
		if index == len(p.plan.Steps)-1 {
			break
		}
	}
	return nil
}

func (p *AgentService) getStepInfo() (int, string, error) {
	steps := p.plan.Steps

	for i, step := range steps {
		re := regexp.MustCompile(`\[([A-Z_]+)\]`)
		match := re.FindStringSubmatch(step.Content)
		if len(match) != 0 {
			typeMath := match[1]
			fmt.Printf("step %d type %s", i, typeMath)
		}

		if step.State == NO_STARTED {
			err := p.markStep(i, IN_PROGRESS)
			if err != nil {
				return 0, "", err
			}
			return i, step.Content, nil
		}
	}
	return -1, "", nil
}

func (p *AgentService) executeStep(ctx context.Context, executor *PlanExecutor, index int, step string) error {
	plan := p.formatPlan()
	stepPrompt := fmt.Sprintf(`
CURRENT PLAN STATUS:
%s
YOUR CURRENT TASK:
You are now working on step %d: %s
Please execute this step using the appropriate tools. When you're done, provide a summary of what you accomplished.
`, plan, index, step)

	str, err := executor.Run(ctx, stepPrompt)
	if err != nil {
		return err
	}

	fmt.Println(str)
	err = p.markStep(index, COMPLETED)

	if err != nil {
		return err
	}
	return nil
}

func (p *AgentService) newPlanTool() domain.Tool {
	var t domain.Tool
	t.Type = "function"
	t.Function = domain.Function{
		Name:        "planning",
		Description: "A planning that allows the agent to create and manage plans for solving complex tasks.\n The provides functionality for creating plans, updating plan steps, and tracking progress.",
		Parameters: &domain.FunctionParameters{
			Properties: params.NewPlanParams(),
			Required:   []string{"command"},
		},
	}

	return t
}

func (p *AgentService) initPlanWithArgs(args string) error {
	var (
		parsedArgs map[string]interface{}
		err        error
	)

	if err = json.Unmarshal([]byte(args), &parsedArgs); err != nil {
		return err
	}

	cmd := parsedArgs["command"].(string)
	if cmd != "create" {
		return errors.New("args 非 create")
	}

	p.plan.Title = parsedArgs["title"].(string)

	steps := parsedArgs["steps"].([]interface{})
	p.plan.Steps = make([]domain.Step, len(steps))
	for i, step := range steps {
		if s, ok := step.(string); ok {
			p.plan.Steps[i] = domain.Step{State: NO_STARTED, Content: s}
		}
	}
	return nil
}

func (p *AgentService) markStep(index int, state string) error {
	if index >= len(p.plan.Steps) {
		return errors.New("当前 step index 非法")
	}
	p.plan.Steps[index].State = state
	return nil
}

func (p *AgentService) formatPlan() string {
	output := fmt.Sprintf("Plan: %s (ID: %s)\n", p.plan.Title, p.Id)
	output += strings.Repeat("=", len(output)) + "\n\n"
	total := len(p.plan.Steps)

	noStarted := 0
	progress := 0
	completed := 0
	blocked := 0

	for _, step := range p.plan.Steps {
		if step.State == NO_STARTED {
			noStarted += 1
		}

		if step.State == IN_PROGRESS {
			progress += 1
		}

		if step.State == COMPLETED {
			completed += 1
		}

		if step.State == BLOCKED {
			blocked += 1
		}
	}

	output += fmt.Sprintf("Progress: %d / %d steps completed ", completed, total)
	if total > 0 {
		percentage := float64(completed) / float64(total) * 100
		output += fmt.Sprintf("(%.1f%%)\n", percentage)
	} else {
		output += "(0%)\n"
	}

	output += fmt.Sprintf("status: %d completed, %d progress, %d no strated, %d blocked", completed, progress, noStarted, blocked)
	output += "Steps:\n"

	statusSymbol := "[ ]"
	for i, step := range p.plan.Steps {
		switch step.State {
		case NO_STARTED:
			statusSymbol = "[ ]"
		case IN_PROGRESS:
			statusSymbol = "[→]"
		case COMPLETED:
			statusSymbol = "[✓]"
		case BLOCKED:
			statusSymbol = "[!]"
		}
		output += fmt.Sprintf("%d. %s %s\n", i, statusSymbol, step.Content)
	}
	return output
}
