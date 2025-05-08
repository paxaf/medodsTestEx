package httpserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/medodsTestEx/internal/usecase"
)

type FeedbackHandler struct {
	usecase usecase.UseCase
}

type errorResponse struct {
	Error string `json:"erros"`
}

func NewFeedbackHandler(usecase usecase.UseCase) *FeedbackHandler {
	return &FeedbackHandler{
		usecase: usecase,
	}
}

func (h *FeedbackHandler) SubmitPing(c *gin.Context) {
	c.String(200, "pong")
}

func (h *FeedbackHandler) GetTokens(c *gin.Context) {
	guid := c.Query("guid")
	tokens, err := h.usecase.GetTokens(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
		return
	}
	log.Println(tokens)
	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		60*60*24*7,
	)
	c.JSON(200, tokens)
}
