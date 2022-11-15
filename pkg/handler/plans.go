package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createGroupPlans(c *gin.Context) {
	ByteBody, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(ByteBody))

	err := h.services.Plans.CreatePlans("", ByteBody)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) GetMasPlans(c *gin.Context) {
	guid_program := c.Params.ByName("guid_program")
	mas, err := h.services.GetPlans(guid_program)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//fmt.Println(mas)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	c.JSON(http.StatusOK, mas)
}

func (h *Handler) GetWorkProgram(c *gin.Context) {
	guid_plan := c.Params.ByName("guid_plan")
	workProgram, err := h.services.GetWorkProgram(guid_plan)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	c.JSON(http.StatusOK, workProgram)
}

func (h *Handler) SavePlan(c *gin.Context) {
	guid_plan := c.Params.ByName("guid_plan")
	key_field := c.Params.ByName("key_field")
	type RequestTask struct {
		Text string `json:"text"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.services.Plans.SavePlan(guid_plan, key_field, rt.Text)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetField(c *gin.Context) {
	guid_plan := c.Params.ByName("guid_plan")
	key_field := c.Params.ByName("key_field")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	field, err := h.services.Plans.GetField(guid_plan, key_field)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, field)
}
