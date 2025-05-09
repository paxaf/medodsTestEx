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

// @Summary Получение токенов
// @Description Генерирует и возвращает JWT и refresh-токен для указанного GUID
// @Tags auth
// @Accept json
// @Produce json
// @Param guid query string true "Уникальный идентификатор пользователя"
// @Success 200 {object} tokens.Tokens
// @Failure 400 {object} map[string]string "Ошибка запроса"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /tokens [get]
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

// @Summary Валидация токенов и получение GUID
// @Description Проверяет валидность токенов и возвращает GUID пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param tokens body tokens.Tokens true "Токены для валидации"
// @Success 200 {object} map[string]string "GUID пользователя"
// @Failure 400 {object} map[string]string "Ошибка запроса"
// @Failure 401 {object} map[string]string "Невалидные токены"
// @Router /guid [get]
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

// @Summary Обновление токенов
// @Description Обновляет JWT и refresh-токен, если текущие валидны
// @Tags auth
// @Accept json
// @Produce json
// @Param tokens body tokens.Tokens true "Текущие токены"
// @Success 200 {object} tokens.Tokens "Новые токены"
// @Failure 400 {object} map[string]string "Ошибка запроса"
// @Failure 401 {object} map[string]string "Невалидные токены"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /refresh [post]
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

// @Summary Деавторизация пользователя
// @Description Деавторизует пользователя по его GUID (инвалидирует токены)
// @Tags auth
// @Accept json
// @Produce json
// @Param tokens body tokens.Tokens true "Токены для деавторизации"
// @Success 200 {object} map[string]string "Успешная деавторизация"
// @Failure 400 {object} map[string]string "Ошибка запроса"
// @Failure 401 {object} map[string]string "Невалидные токены"
// @Router /deauth [get]
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
