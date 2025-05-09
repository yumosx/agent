package domain

import (
	"github.com/yumosx/agent/internal/domain/params"
)

type LLMRequest struct {
	SystemContent string
	Msgs          []Msg
	Tools         []Tool
	Choice        string
}

var (
	SYSTEM    = "SYSTEM"
	USER      = "user"
	ASSISTANT = "assistant"
	TOOL      = "tool"
)

type Msg struct {
	Role      string
	Content   string
	Id        string
	ToolCalls []LLMToolCall
}

type Tool struct {
	Type     string
	Function Function
}

type Function struct {
	Name        string
	Description string
	Parameters  *FunctionParameters
}

type FunctionParameters struct {
	Type       string
	Properties *params.Parameters
	Required   []string
}

type LLMResponse struct {
	Content   string
	Done      bool
	ToolCalls []LLMToolCall
}

type LLMToolCall struct {
	Index    int
	ID       string
	Type     string
	Function LLMToolCallFunction
}

type LLMToolCallFunction struct {
	Name      string
	Arguments string
}
