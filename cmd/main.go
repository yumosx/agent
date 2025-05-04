package main

import (
	"github.com/cohesion-org/deepseek-go"
	"github.com/gin-gonic/gin"
	"github.com/yumosx/agent/internal/handler"
	"github.com/yumosx/agent/internal/service"
	"github.com/yumosx/agent/internal/service/llm"
	"os"
)

func main() {
	token := os.Getenv("token")
	client := deepseek.NewClient(token)
	llmHandler := llm.NewHandler(client)
	executor := service.NewPlanExecutor(llmHandler)
	plan := service.NewPlanService("1", llmHandler, executor)
	hd := handler.NewHandler(plan)

	router := gin.Default()
	hd.SetupRoutes(router)
	router.Run(":8080")
}
