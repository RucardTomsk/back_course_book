package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccessInput struct {
	GuidUser string `json:"guiduser"`
	GuidNode string `json:"guidnode"`
}

func (h *Handler) IssueAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	var ai AccessInput
	if err := c.BindJSON(&ai); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user_guid, _ := c.Get(userGuid)
	flag_admin, err := h.services.CheckRoleAdmin(user_guid.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	flag_access, err := h.services.CheckAccess(user_guid.(string), ai.GuidNode)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !flag_admin && !flag_access {
		newErrorResponse(c, http.StatusForbidden, "not access")
		return
	}

	label, err := h.services.IssueAccess(ai.GuidUser, ai.GuidNode)
	if err != nil {
		fmt.Println(1)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if label == "Faculty" {
		mas_program, err := h.services.GetMasProgram(ai.GuidNode)
		if err != nil {
			fmt.Println(2)
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		for _, guid_program := range mas_program {

			mas_plan, err := h.services.GetPlans(guid_program.Guid)
			if err != nil {
				fmt.Println(3)
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			for _, value := range mas_plan {
				for _, brefPlan := range value {
					_, err := h.services.IssueAccess(ai.GuidUser, brefPlan.Guid)
					if err != nil {
						fmt.Println(4)
						newErrorResponse(c, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}

	}

	c.JSON(http.StatusOK, nil)
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
	GuidNode string
}

func (h *Handler) CheackAccess(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)
	var s_node CheackNode
	if err := c.BindJSON(&s_node); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	flag, err := h.services.CheckAccess(user_guid.(string), s_node.GuidNode)
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
