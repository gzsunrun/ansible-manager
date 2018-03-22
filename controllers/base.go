package controllers

import (
	"github.com/astaxie/beego"
)

// ErrorMsg the struct of error message
type ErrorMsg struct {
	Code    int    `json:"err_code"`
	Message string `json:"message"`
}

// BaseController the base controller
type BaseController struct {
	beego.Controller
}

// GetUid get user uuid
func (c *BaseController) GetUid() string {
	return c.Ctx.Input.GetData("uid").(string)
}

// SetResult return json to http
func (c *BaseController) SetResult(err error, result interface{}, errcode int, key ...string) {
	c.Ctx.Output.Status = errcode
	if err != nil {
		c.Data["json"] = result
		return
	}
	if result == nil && (len(key) == 0 || key[0] == "") {
		return
	}

	if len(key) == 0 || key[0] == "" {
		c.Data["json"] = result
	} else {
		c.Data["json"] = map[string]interface{}{key[0]: result}
	}
}
