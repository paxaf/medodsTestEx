package httpserver

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/medodsTestEx/internal/tokens"
	"github.com/paxaf/medodsTestEx/internal/usecase"
)

type FeedbackHandler struct {
	usecase usecase.UseCase
}

func hashUserAgent(ua string) string {
	h := sha256.Sum256([]byte(ua))
	return hex.EncodeToString(h[:])
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
	ua := c.GetHeader("User-Agent")
	hashUa := hashUserAgent(ua)
	guid := c.Query("guid")
	tokens, err := h.usecase.GetTokens(guid, hashUa)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, tokens)
}

func (h *FeedbackHandler) Guid(c *gin.Context) {
	var req tokens.Tokens
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	log.Println(req.RefreshToken)
	guid, _, valid := h.usecase.ValidateTokens(req)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthtorized"})
		return
	}
	c.JSON(200, map[string]string{"guid": guid})
}

func (h *FeedbackHandler) Refresh(c *gin.Context) {
	var req tokens.Tokens
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	guid, hashedAgent, valid := h.usecase.ValidateTokens(req)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthtorized"})
		return
	}
	ua := c.GetHeader("User-Agent")
	hashUa := hashUserAgent(ua)
	if hashUa != hashedAgent {
		h.usecase.UnathorizeUser(guid)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "rejected"})
		return
	}
	token, err := h.usecase.UpdateTokens(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internalError"})
		return
	}
	c.JSON(200, token)
}

func (h *FeedbackHandler) Deauthorized(c *gin.Context) {
	var req tokens.Tokens
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	guid, valid := h.usecase.ValidateJWT(req)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthtorized"})
		return
	}
	h.usecase.UnathorizeUser(guid)
	c.JSON(200, map[string]string{})
}
