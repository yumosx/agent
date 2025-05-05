package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yumosx/agent/internal/service"
	"net/http"
)

type Handler struct {
	svc *service.PlanService
}

func NewHandler(svc *service.PlanService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/", h.serveIndex)
	router.POST("/chat", h.handleChat)
	router.POST("/code", h.handleCode)
}

func (h *Handler) serveIndex(ctx *gin.Context) {
	ctx.File("./internal/font/index.html")
}

func (h *Handler) handleChat(ctx *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	plan, err := h.svc.Plan(ctx, request.Message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "内部错误"})
		return
	}

	response := gin.H{
		"response": plan,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) handleCode(ctx *gin.Context) {

}
