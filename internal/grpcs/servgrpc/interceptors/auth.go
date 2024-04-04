package interceptors

import (
	"context"
	"strings"

	"github.com/shulganew/shear.git/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	zap.S().Infoln("Intersepted!")

	pass := ctx.Value(config.CtxPassKey{}).(string)
	zap.S().Infoln("Pass!", pass)
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["user_id"]) {

		return nil, status.Errorf(codes.Unauthenticated, "user unauthorized")

	}
	m, err := handler(ctx, req)
	if err != nil {
		zap.S().Infoln("RPC failed with error: %v", err)
	}
	// to add context
	//https://stackoverflow.com/questions/71114401/grpc-how-to-pass-value-from-interceptor-to-service-function
	return m, err
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

/*

}


// Middleware function check user's authorization.
//
// Verify encrypted cookie, check or add if not existed context context variable config.CtxConfig user_id(uuid) and new_user(bool).
func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// get password from context
		pass := req.Context().Value(config.CtxPassKey{}).(string)
		isNewUser := false
		userID, ok := service.GetCodedUserID(req, pass)

		if ok {
			// cookie user_id is set
			cookies := req.Cookies()

			// clean cookie data
			req.Header["Cookie"] = make([]string, 0)
			for _, cookie := range cookies {

				if cookie.Name == "user_id" {
					cookie.Value = userID
				}
				req.AddCookie(cookie)
			}
		} else {
			// cookie not set or not decoded
			// create new user uuid
			userID, err := uuid.NewV7()
			if err != nil {
				zap.S().Errorln("Error generate user uuid")
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}

			// encode cookie for client
			coded, err := service.EncodeCookie(userID.String(), pass)
			if err != nil {
				zap.S().Errorln("Error encode uuid")
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			// set to response
			codedCookie := http.Cookie{Name: "user_id", Value: coded}
			http.SetCookie(res, &codedCookie)

			// set to request
			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)

			// mark new user for handlers
			newUser := http.Cookie{Name: "new_user", Value: "true"}
			req.AddCookie(&newUser)
			isNewUser = true

		}
		ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userID, isNewUser))
		h.ServeHTTP(res, req.WithContext(ctx))
	})

}
*/
