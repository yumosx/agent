package service

import (
	"context"
	"github.com/cohesion-org/deepseek-go"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/agent/internal/service/agent"
	"github.com/yumosx/agent/internal/service/llm"
	"os"
	"testing"
)

func TestPlanExecute(t *testing.T) {
	token := os.Getenv("token")
	client := deepseek.NewClient(token)
	handler := llm.NewHandler(client)

	printer := agent.NewPrintExecutor()

	plan := NewPlanService("1", handler, []agent.Executor{printer})
	err := plan.Execute(context.Background(), "帮我生成一个优化 SQL 的性能分析计划")
	require.NoError(t, err)
}
