package service

import (
	"context"
	"fmt"
	"github.com/yumosx/agent/internal/domain"
	"github.com/yumosx/agent/internal/domain/params"
	"github.com/yumosx/agent/internal/service/llm"
)

type PlanExecutor struct {
	maxStep int
	tools   []domain.Tool
	handler *llm.Handler
	// 用户模型的上下文
	messages []domain.Msg
	results  []string
}

const (
	system   = `You are an agent that can execute tool calls"`
	nextStep = `Based on user needs, proactively select the most appropriate tool or combination of tools. For complex tasks, you can break down the problem and use different tools step by step to solve it. After using each tool, clearly explain the execution results and suggest the next steps.
If you want to stop the interaction at any point, use the "terminate" tool/function call.`
)

func NewPlanExecutor(handler *llm.Handler) *PlanExecutor {
	return &PlanExecutor{handler: handler, messages: make([]domain.Msg, 10)}
}

func (p *PlanExecutor) Run(ctx context.Context, step string) (string, error) {
	p.messages = append(p.messages, domain.Msg{Role: domain.USER, Content: step})
	result, err := p.step(ctx)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *PlanExecutor) step(ctx context.Context) (string, error) {
	var req domain.LLMRequest
	req.SystemContent = system

	req.Msgs = make([]domain.Msg, len(p.messages))
	for _, msg := range p.messages {
		req.Msgs = append(req.Msgs, msg)
	}

	req.Msgs = append(req.Msgs, domain.Msg{
		Role:    domain.USER,
		Content: nextStep,
	})
	req.Tools = []domain.Tool{p.newChatTool(), p.newTrimTool(), p.newGoTool()}
	resp, err := p.handler.Invoke(ctx, req)
	if err != nil {
		return "", err
	}

	p.messages = append(p.messages, domain.Msg{Role: domain.ASSISTANT, Content: resp.Content})

	toolCalls := resp.ToolCalls
	var result string

	for _, tool := range toolCalls {
		switch tool.Function.Name {
		case "terminate":
			result = p.executeTrim(tool.Function.Arguments)
		case "golang_execute":
			result = p.executeGolang(tool.Function.Arguments)
		}
		println(result)
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
				Properties: params.NewChatParams(),
				Required:   []string{"response"},
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
				Properties: params.NewTrimParams(),
				Required:   []string{"status"},
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
				Properties: params.NewGoParams(),
				Required:   []string{"code"},
			},
		},
	}
}

func (p *PlanExecutor) executeTrim(str string) string {
	return fmt.Sprintf("The interaction has been completed with status: %s", str)
}

func (p *PlanExecutor) executeGolang(str string) string {
	return fmt.Sprintf("hello")
}
