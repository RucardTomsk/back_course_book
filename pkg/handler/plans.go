package handler

import (
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createGroupPlans(c *gin.Context) {
	guid_faculty := c.Params.ByName("guid_faculty")

	formFile, _ := c.FormFile("file")
	openedFile, _ := formFile.Open()
	file, _ := ioutil.ReadAll(openedFile)

	err := h.services.Plans.CreatePlans("", file, guid_faculty)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, mas)
}

func (h *Handler) GetWorkProgram(c *gin.Context) {
	guid_plan := c.Params.ByName("guid_plan")
	workProgram, err := h.services.GetWorkProgram(guid_plan)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.services.Plans.SavePlan(guid_plan, key_field, rt.Text)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetField(c *gin.Context) {
	guid_plan := c.Params.ByName("guid_plan")
	key_field := c.Params.ByName("key_field")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	field, err := h.services.Plans.GetField(guid_plan, key_field)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, field)
}

func (h *Handler) CopyPlan(c *gin.Context) {
	guid_from := c.Params.ByName("guid_from")
	guid_to := c.Params.ByName("guid_to")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	if err := h.services.Plans.CloneFieldPlan(guid_from, guid_to); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) GenerateWord(c *gin.Context) {
	formFile, _ := c.FormFile("file")
	openedFile, _ := formFile.Open()
	file, _ := ioutil.ReadAll(openedFile)
	var permissions fs.FileMode
	permissions = 0644 // or whatever you need
	err := os.WriteFile("file.pdf", file, permissions)
	if err != nil {
		// handle error
	}
	//res,,err := docconv.ConvertHTML()

}

func (h *Handler) GetNamePlans(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	guid_plan := c.Params.ByName("guid")
	mas, err := h.services.GetNamePlans(guid_plan)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, mas)
}
