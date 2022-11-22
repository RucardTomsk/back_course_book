package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IssueAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)
	flag_admin, err := h.services.CheckRoleAdmin(user_guid.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	type accessInput struct {
		guidUser string `json:"g_u"`
		guidNode string `json:"g_n"`
	}

	var inputData accessInput
	if err := c.BindJSON(&inputData); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	flag_access, err := h.services.CheckAccess(user_guid.(string), inputData.guidNode)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !flag_admin || !flag_access {
		newErrorResponse(c, http.StatusForbidden, "not access")
		return
	}

	label, err := h.services.IssueAccess(inputData.guidUser, inputData.guidNode)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if label == "Faculty" {
		mas_program, err := h.services.GetMasProgram(inputData.guidNode)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		for _, guid_program := range mas_program {

			mas_plan, err := h.services.GetPlans(guid_program.Guid)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			for _, value := range mas_plan {
				for _, brefPlan := range value {
					_, err := h.services.IssueAccess(inputData.guidUser, brefPlan.Guid)
					if err != nil {
						newErrorResponse(c, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}

	}

	c.JSON(http.StatusOK, inputData)
}

func (h *Handler) CheackAdmin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)
	flag, err := h.services.CheckRoleAdmin(user_guid.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !flag {
		newErrorResponse(c, http.StatusForbidden, "not access")
		return
	}

	c.JSON(http.StatusOK, nil)
}

type CheackNode struct {
	guid_node string
}

func (h *Handler) CheackAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)
	var s_node CheackNode
	if err := c.BindJSON(&s_node); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	flag, err := h.services.CheckAccess(user_guid.(string), s_node.guid_node)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !flag {
		newErrorResponse(c, http.StatusForbidden, "not access")
		return
	}

	c.JSON(http.StatusOK, nil)
}
