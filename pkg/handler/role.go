package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
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

type SInvute struct {
	GuidNode string
	Email    string
}

func (h *Handler) createInvite(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)
	var s_invite SInvute
	if err := c.BindJSON(&s_invite); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Print(s_invite.Email, s_invite.GuidNode)
	flag_admin, err := h.services.CheckRoleAdmin(user_guid.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	flag, err := h.services.CheckAccess(user_guid.(string), s_invite.GuidNode)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !flag_admin && !flag {
		newErrorResponse(c, http.StatusForbidden, "not access")
		return
	}

	fmt.Println("guid_node", s_invite.GuidNode)

	guid_invite, err := h.services.CreateInvite(s_invite.GuidNode)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	name_node, err := h.services.Plans.GetField(s_invite.GuidNode, "name")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("name_node", name_node)
	fmt.Println("email", s_invite.Email)

	msg := gomail.NewMessage()
	msg.SetHeader("From", "course-book@tsu.ru")
	msg.SetHeader("To", s_invite.Email)
	msg.SetHeader("Subject", "Приглашение на редактирование дисциплины")
	msg.SetBody("text/html", fmt.Sprintf(`<p>Вы были приглашены как преподаватель/редактор на дисциплину "%s"</p>
										<p>Ваша ссылка-приглашение</p>
										<a href="http://localhost:8080/#/invite/%s">http://localhost:8080/#/invite/%s</a>`, name_node, guid_invite, guid_invite))
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	n := gomail.NewDialer("smtp.gmail.com", 465, "www.carat.ru@gmail.com", os.Getenv("PASS_EMAIL"))

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nil)
}

type CheackInvite struct {
	GuidInvite string
}

func (h *Handler) useInvite(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	user_guid, _ := c.Get(userGuid)

	var s_node CheackInvite
	if err := c.BindJSON(&s_node); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UseInvite(s_node.GuidInvite, user_guid.(string)); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
