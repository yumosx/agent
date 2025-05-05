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
	plan := NewPlanService(handler, executor)
	var ctx = context.Background()

	s, err := plan.Plan(ctx, "执行ls 命令")
	require.NoError(t, err)
	println(s)

	require.NoError(t, err)
	err = plan.Execute(context.Background())
	require.NoError(t, err)
}
