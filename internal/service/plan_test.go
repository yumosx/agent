package service

import (
	"github.com/cohesion-org/deepseek-go"
	"github.com/yumosx/agent/internal/service/llm"
	"os"
	"testing"
)

func TestPlanExecute(t *testing.T) {
	token := os.Getenv("token")
	client := deepseek.NewClient(token)
	handler := llm.NewHandler(client)

	executor := NewPlanExecutor(handler)
	plan := NewPlanService("1", handler, executor)
	_ = plan
}
