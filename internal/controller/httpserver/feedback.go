package httpserver

import (
	"log"

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
		log.Println(err)
		return
	}
	c.JSON(200, tokens)
}

func (h *FeedbackHandler) Guid(c *gin.Context) {

}

func (h *FeedbackHandler) Refresh(c *gin.Context) {

}
