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

	if len(req.Msgs) != 0 {
		for _, msg := range req.Msgs {
			if msg.Role == domain.USER {
				request.Messages = append(request.Messages,
					deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleUser, Content: msg.Content})
			}
			if msg.Role == domain.ASSISTANT {
				request.Messages = append(request.Messages,
					deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleAssistant, Content: msg.Content})
			}
			if msg.Role == domain.FUNCTION {
				request.Messages = append(request.Messages,
					deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleTool, Content: msg.Content, ToolCallID: msg.Id})
			}
		}
	}

	if len(req.Tools) != 0 {
		request.Tools = make([]deepseek.Tool, len(req.Tools))
		for i, tool := range req.Tools {
			request.Tools[i].Type = tool.Type
			request.Tools[i].Function.Name = tool.Function.Name
			request.Tools[i].Function.Description = tool.Function.Description
			request.Tools[i].Function.Parameters = &deepseek.FunctionParameters{
				Type:       "object",
				Properties: tool.Function.Parameters.Properties.ToMap(),
				Required:   tool.Function.Parameters.Required,
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
