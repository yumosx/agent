package service

import (
	"context"
	"github.com/yumosx/agent/internal/domain"
	"github.com/yumosx/agent/internal/service/llm"
)

type PlanExecutor struct {
	maxStep int
	tools   []domain.Tool
	handler *llm.Handler
	// 用户模型的上下文
	messages []string
	results  []string
}

const (
	system   = `You are an agent that can execute tool calls"`
	nextStep = ` Based on user needs, proactively select the most appropriate tool or combination of tools. For complex tasks, you can break down the problem and use different tools step by step to solve it. After using each tool, clearly explain the execution results and suggest the next steps.
If you want to stop the interaction at any point, use the "terminate" tool/function call.`
)

func NewPlanExecutor(handler *llm.Handler) *PlanExecutor {
	return &PlanExecutor{handler: handler, messages: make([]string, 10)}
}

func (p *PlanExecutor) GetType() string {
	return "printer"
}

func (p *PlanExecutor) Run(step string) (string, error) {
	return p.Step(step)
}

func (p *PlanExecutor) Step(step string) (string, error) {
	p.messages = append(p.messages, step)
	result, err := p.step()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *PlanExecutor) step() (string, error) {
	var req domain.LLMRequest
	req.SystemContent = system

	req.Content = ""
	for _, msg := range p.messages {
		req.Content += msg
	}

	req.Content += nextStep
	req.Tools = []domain.Tool{p.newChatTool(), p.newTrimTool(), p.newGoTool()}
	resp, err := p.handler.Invoke(context.Background(), req)
	if err != nil {
		return "", err
	}

	toolCalls := resp.ToolCalls
	var result string

	for _, tool := range toolCalls {
		switch tool.Function.Name {
		case "terminate":
			result = p.executeTrim(tool.Function.Arguments)
		case "golang_execute":
			result = p.executeGolang(tool.Function.Arguments)
		}
	}

	return result, nil
}

func (p *PlanExecutor) newChatTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name:        "create_chat_completion",
			Description: "Creates a structured completion with specified output formatting.",
			Parameters: &domain.FunctionParameters{
				Required: []string{"response"},
			},
		},
	}
}

func (p *PlanExecutor) newTrimTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name: "terminate",
			Description: `Terminate the interaction when the request is met OR if the assistant cannot proceed further with the task.  
When you have finished all the tasks, call this tool to end the work.`,
			Parameters: &domain.FunctionParameters{
				Required: []string{"status"},
			},
		},
	}
}

func (p *PlanExecutor) newGoTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name: "golang_execute",
			Description: `Executes Golang code string. Note: Only print outputs are visible, function return values are not captured. 
Use print statements to see results.`,
			Parameters: &domain.FunctionParameters{
				Required: []string{"code"},
			},
		},
	}
}

func (p *PlanExecutor) executeTrim(str string) string {
	return ""
}

func (p *PlanExecutor) executeGolang(str string) string {
	return ""
}
