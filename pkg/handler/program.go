package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMasProgram(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	guid_faculty := c.Params.ByName("guid_faculty")
	mas_program, err := h.services.GetMasProgram(guid_faculty)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, mas_program)

}
