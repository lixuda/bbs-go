package middleware

import (
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/cache"
	"bbs-go/model"
)

// AdminAuth 后台权限
func AdminAuth(ctx iris.Context) {
	token := getUserToken(ctx)
	userToken := cache.UserTokenCache.Get(token)

	// 没找到授权
	if userToken == nil || userToken.Status == model.StatusDeleted {
		notLogin(ctx)
		return
	}
	// 授权过期
	if userToken.ExpiredAt <= simple.NowTimestamp() {
		notLogin(ctx)
		return
	}

	user := cache.UserCache.Get(userToken.UserId)
	if user == nil || !services.UserService.HasRole(user, model.RoleOwner) {
		_, _ = ctx.JSON(simple.JsonErrorCode(2, "无权限"))
		ctx.StopExecution()
		return
	}

	ctx.Next()
}

// 从请求体中获取UserToken
func getUserToken(ctx iris.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

func notLogin(ctx iris.Context) {
	_, _ = ctx.JSON(simple.JsonError(simple.ErrorNotLogin))
	ctx.StopExecution()
}
