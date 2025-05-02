package domain

type LLMResponse struct {
	Content   string
	Done      bool
	ToolCalls []LLMToolCall
}

type LLMRequest struct {
	SystemContent string
	Content       string
	Tools         []Tool
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
	Properties map[string]interface{}
	Required   []string
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
