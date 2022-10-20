package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thinkeridea/go-extend/exnet"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
)

// 传递公共参数中间件
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//为了测试loading状态，生产版需要删掉
		time.Sleep(time.Duration(300) * time.Millisecond)

		reqToken := r.Header.Get(consts.AUTHORIZATION)
		splitToken := strings.Split(reqToken, consts.BEARER)
		v := contexts.ContextValues{}
		if len(splitToken) == 2 {
			reqToken = splitToken[1]
			if reqToken != "" {
				v.Token = reqToken
				me, err := authentication.GetUserByToken(reqToken)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				v.Me = me
			}
		}
		appId := r.Header.Get(consts.HEADER_APPX_APPID)
		if appId == "" {
			v.AppId = app.GetSystemApp().AppId
		} else {
			intAppId, _ := strconv.ParseUint(appId, 10, 64)
			v.AppId = intAppId
		}

		v.Host = r.Host
		ip := exnet.ClientPublicIP(r)
		if ip == "" {
			ip = exnet.ClientIP(r)
		}
		v.IP = ip
		ctx := context.WithValue(r.Context(), consts.CONTEXT_VALUES, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}