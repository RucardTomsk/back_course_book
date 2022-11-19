package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMasFaculty(c *gin.Context) {
	faculty, err := h.services.GetMasFaculty()
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, faculty)

}

func (h *Handler) GetNameFaculte(c *gin.Context) {
	guid := c.Params.ByName("guid")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	namef, err := h.services.GetNameFaculty(guid)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, namef)
}
