package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yumosx/agent/internal/domain"
	"github.com/yumosx/agent/internal/domain/params"
	"github.com/yumosx/agent/internal/service/llm"
	"github.com/yumosx/agent/internal/tool"
	"time"
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
	system   = `You are an agent that can execute tool calls, only call tools when they are absolutely necessary. if the user's task is general or you already know the answer, respond without calling tools."`
	nextStep = `Based on user needs, proactively select the most appropriate tool or combination of tools. For complex tasks, you can break down the problem and use different tools step by step to solve it. After using each tool, clearly explain the execution results and suggest the next steps.
If you want to stop the interaction at any point, use the "terminate" tool/function call.`
)

func NewPlanExecutor(handler *llm.Handler) *PlanExecutor {
	return &PlanExecutor{handler: handler, messages: make([]domain.Msg, 0)}
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

	req.Msgs = make([]domain.Msg, 0)
	for _, msg := range p.messages {
		req.Msgs = append(req.Msgs, msg)
	}

	req.Msgs = append(req.Msgs, domain.Msg{
		Role:    domain.USER,
		Content: nextStep,
	})
	req.Tools = []domain.Tool{p.newChatTool(), p.newTrimTool(), p.newGoTool(), p.newBashTool()}
	resp, err := p.handler.Invoke(ctx, req)
	if err != nil {
		return "", err
	}

	p.messages = append(p.messages, domain.Msg{Role: domain.ASSISTANT, Content: resp.Content})

	toolCalls := resp.ToolCalls
	var result string

	for _, t := range toolCalls {
		switch t.Function.Name {
		case "terminate":
			return t.Function.Arguments, nil
		case "golang_execute":
			p.executeCode(t.Function.Arguments)
		case "bash":
			p.executeBash(t.Function.Arguments)
		default:
			continue
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

func (p *PlanExecutor) newBashTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name: "bash",
			Description: `Execute a bash command in the terminal.
* Long running commands: For commands that may run indefinitely, it should be run in the background and the output should be redirected to a file, e.g. command = "python3 app.py > server.log 2>&1 &".
* Interactive: If a bash command returns exit code "-1", this means the process is not yet finished. The assistant must then send a second call to terminal with an empty "command" (which will retrieve any additional logs), or it can send additional text (set "command" to the text) to STDIN of the running process, or it can send command="ctrl+c" to interrupt the process.
* Timeout: If a command execution result says "Command timed out. Sending SIGINT to the process", the assistant should retry running the command in the background.`,
			Parameters: &domain.FunctionParameters{
				Properties: params.NewBashParams(),
				Required:   []string{"command"},
			},
		},
	}
}

func (p *PlanExecutor) newGoTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name:        "golang_execute",
			Description: `Executes Golang code string. Note: Only print outputs are visible, function return values are not captured. Use print statements to see results.`,
			Parameters: &domain.FunctionParameters{
				Properties: params.NewGoParams(),
				Required:   []string{"code"},
			},
		},
	}
}

func (p *PlanExecutor) newGolangTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name:        "coder",
			Description: ``,
			Parameters: &domain.FunctionParameters{
				Properties: params.NewPlanParams(),
				Required:   []string{""},
			},
		},
	}
}

func (p *PlanExecutor) newBrowserTool() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name:        "browser_use",
			Description: ``,
			Parameters: &domain.FunctionParameters{
				Properties: params.NewBrowserUse(),
				Required:   []string{""},
			},
		},
	}
}

func (p *PlanExecutor) newFileWrite() domain.Tool {
	return domain.Tool{
		Type: "function",
		Function: domain.Function{
			Name:        "",
			Description: "",
			Parameters: &domain.FunctionParameters{
				Properties: nil,
				Required:   []string{""},
			},
		},
	}
}

func (p *PlanExecutor) executeTrim(status string) string {
	return fmt.Sprintf("The interaction has been completed with status: %s", status)
}

func (p *PlanExecutor) executeBash(args string) string {
	var cmd map[string]string
	if err := json.Unmarshal([]byte(args), &cmd); err != nil {
		return fmt.Sprintf("response format umarshal failed: %s", err.Error())
	}

	c := cmd["command"]

	bash := tool.NewBashSession(20 * time.Second)
	err := bash.Start()
	if err != nil {
		return err.Error()
	}

	result, errOutput, err := bash.Run(c)

	if err != nil {
		return err.Error()
	}

	if errOutput != "" {
		return errOutput
	}

	return result
}

func (p *PlanExecutor) executeCode(args string) string {
	var code map[string]string
	if err := json.Unmarshal([]byte(args), &code); err != nil {
		return fmt.Sprintf("response format umarshal failed: %s", err.Error())
	}

	c := code["code"]
	println(c)

	return c
}
