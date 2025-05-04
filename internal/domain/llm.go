package domain

import (
	"github.com/yumosx/agent/internal/domain/params"
)

type LLMRequest struct {
	SystemContent string
	Content       string
	Assistant     string
	Tools         []Tool
	Choice        string
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
