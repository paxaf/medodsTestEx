package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/paxaf/medodsTestEx/internal/usecase"
)

type FeedbackHandler struct {
}

type errorResponse struct {
	Error string `json:"erros"`
}

func NewFeedbackHandler() *FeedbackHandler {
	return &FeedbackHandler{}
}

func (h *FeedbackHandler) SubmitPing(c *gin.Context) {
	c.String(200, "pong")
}

func (h *FeedbackHandler) GetTokens(c *gin.Context) {
	guid := c.Query("guid")
	accessToken, err := usecase.GetTokens(guid)
	if err != nil {
		return
	}
	refreshToken, err :=
		c.String(200, token)
}
