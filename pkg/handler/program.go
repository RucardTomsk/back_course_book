package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMasProgram(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	guid_faculty := c.Params.ByName("guid_faculty")
	mas_program, err := h.services.GetMasProgram(guid_faculty)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, mas_program)

}

func (h *Handler) GetNameProgramAndFaculty(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	guid_program := c.Params.ByName("guid")
	mas_s, err := h.services.GetNameProgramAndFaculty(guid_program)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, mas_s)
}
