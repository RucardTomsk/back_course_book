package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RucardTomsk/course_book/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) register(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) login(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, refreshToken, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) test(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	id, _ := c.Get(userGuid)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetUserFIO(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	guid, _ := c.Get(userGuid)
	fio, err := h.services.GetUserFioByGuid(guid.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fio)
}

func (h *Handler) GetUserNotAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	guid_node := c.Params.ByName("guid")
	mas_user, err := h.services.GetUserNotAccess(guid_node)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(mas_user)
	c.JSON(http.StatusOK, mas_user)
}

type RefreshStruct struct {
	refreshToken string
}

func (h *Handler) CheckToken(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	var RefreshS RefreshStruct
	if err := c.BindJSON(&RefreshS); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		user, err := h.services.Authorization.GetUserToRefreshToken(RefreshS.refreshToken)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		token, refreshToken, err := h.services.Authorization.GenerateToken(user.Email, user.Password)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"token":        token,
			"refreshToken": refreshToken,
		})

		return
	}

	c.JSON(http.StatusOK, nil)
}
