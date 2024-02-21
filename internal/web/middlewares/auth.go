package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// get password from context
		pass := req.Context().Value(config.CtxPassKey{}).(string)

		isNewUser := false

		userID, ok := service.GetCodedUserID(req, pass)

		if ok {
			//cookie user_id is set
			cookies := req.Cookies()

			//clean cookie data
			req.Header["Cookie"] = make([]string, 0)
			for _, cookie := range cookies {

				if cookie.Name == "user_id" {
					cookie.Value = userID
				}
				req.AddCookie(cookie)

			}
		} else {
			//cookie not set or not decoded
			//create new user uuid
			userID, err := uuid.NewV7()
			if err != nil {
				zap.S().Errorln("Error generate user uuid")
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}

			//encode coockie for client
			coded, err := service.EncodeCookie(userID.String(), pass)
			if err != nil {
				zap.S().Errorln("Error encode uuid")
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			//set to response
			codedCookie := http.Cookie{Name: "user_id", Value: coded}
			http.SetCookie(res, &codedCookie)

			//set to request
			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)
			//mark new user for handlers
			newUser := http.Cookie{Name: "new_user", Value: "true"}
			req.AddCookie(&newUser)

			isNewUser = true

		}

		ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userID, isNewUser))
		h.ServeHTTP(res, req.WithContext(ctx))

	})

}
