package handler

import (
	"fmt"
	"net/http"

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

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
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

type sGuidUser struct {
	guid_user string
}

func (h *Handler) GetUserNotAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var NewsGuidUser sGuidUser
	if err := c.BindJSON(&NewsGuidUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	mas_user, err := h.services.GetUserNotAccess(NewsGuidUser.guid_user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(mas_user)
	c.JSON(http.StatusOK, mas_user)
}
