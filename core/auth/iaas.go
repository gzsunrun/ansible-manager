package auth
import (
	"github.com/astaxie/beego/context"
)
func IaaSAuthFilter(ctx *context.Context) {
		uid := "iaas-uuid"
		ctx.Input.SetData("uid", uid)
}