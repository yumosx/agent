package service

import (
	"context"
	"github.com/cohesion-org/deepseek-go"
	"github.com/stretchr/testify/require"
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
	err := plan.Execute(context.Background(), "使用 Golang 编写一段 打印 hello word 代码")
	require.NoError(t, err)
}
