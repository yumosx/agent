package llm

import (
	"context"
	"github.com/cohesion-org/deepseek-go"
	"github.com/yumosx/agent/internal/domain"
)

type Handler struct {
	client *deepseek.Client
}

func NewHandler(client *deepseek.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Invoke(ctx context.Context, req domain.LLMRequest) (domain.LLMResponse, error) {
	request := &deepseek.ChatCompletionRequest{
		Model:      deepseek.DeepSeekChat,
		Messages:   []deepseek.ChatCompletionMessage{},
		ToolChoice: "auto",
	}

	if req.SystemContent != "" {
		request.Messages = append(request.Messages,
			deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleSystem, Content: req.SystemContent})
	}

	if req.Content != "" {
		request.Messages = append(request.Messages,
			deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleUser, Content: req.Content})
	}

	if len(req.Tools) != 0 {
		request.Tools = make([]deepseek.Tool, len(req.Tools))

		if len(req.Tools) != 0 {
			request.Tools = make([]deepseek.Tool, len(req.Tools))
			for i, tool := range req.Tools {
				request.Tools[i].Type = tool.Type
				request.Tools[i].Function.Name = tool.Function.Name
				request.Tools[i].Function.Description = tool.Function.Description
				request.Tools[i].Function.Parameters = &deepseek.FunctionParameters{
					Type:     "object",
					Required: tool.Function.Parameters.Required,
				}

				switch tool.Function.Name {
				case "planning":
					request.Tools[i].Function.Parameters.Properties = h.NewPlanParams()
				default:
				}
			}
		}
	}

	response, err := h.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return domain.LLMResponse{}, err
	}

	ch := response.Choices[0].Message

	resp := domain.LLMResponse{Content: ch.Content}

	if len(ch.ToolCalls) == 0 {
		return resp, nil
	}

	resp.ToolCalls = make([]domain.LLMToolCall, 0, len(ch.ToolCalls))

	for _, tool := range ch.ToolCalls {
		resp.ToolCalls = append(resp.ToolCalls, domain.LLMToolCall{
			ID:    tool.ID,
			Index: tool.Index,
			Type:  tool.Type,
			Function: domain.LLMToolCallFunction{
				Name:      tool.Function.Name,
				Arguments: tool.Function.Arguments,
			},
		})
	}

	return resp, nil
}

func (h *Handler) NewPlanParams() map[string]interface{} {
	return map[string]interface{}{
		"command": map[string]interface{}{
			"description": "The command to execute. Available commands: create, update, list, get, set_active, mark_step, delete.",
			"enum": []string{
				"create",
				"update",
				"list",
				"get",
				"set_active",
				"mark_step",
				"delete",
			},
			"type": "string",
		},
		"plan_id": map[string]interface{}{
			"description": "Unique identifier for the plan. Required for create, update, set_active, and delete commands. Optional for get and mark_step (uses active plan if not specified).",
			"type":        "string",
		},
		"title": map[string]interface{}{
			"description": "Title for the plan. Required for create command, optional for update command.",
			"type":        "string",
		},
		"steps": map[string]interface{}{
			"description": "List of plan steps. Required for create command, optional for update command.",
			"type":        "array",
			"items": map[string]string{
				"type": "string",
			},
		},
		"step_index": map[string]interface{}{
			"description": "Index of the step to update (0-based). Required for mark_step command.",
			"type":        "integer",
		},
		"step_status": map[string]interface{}{
			"description": "Status to set for a step. Used with mark_step command.",
			"enum": []string{
				"not_started",
				"in_progress",
				"completed",
				"blocked",
			},
			"type": "string",
		},
		"step_notes": map[string]interface{}{
			"description": "Additional notes for a step. Optional for mark_step command.",
			"type":        "string",
		},
	}
}
