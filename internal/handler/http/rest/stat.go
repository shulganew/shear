package rest

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

// DTO JSON response.
type StatDTO struct {
	// Num total URLs in shortener.
	Shorts int `json:"urls"`
	// Num total URLs in shortener.
	Users int `json:"users"`
}

// Handler for API:
//
//	GET "/api/internal/stats"
type HandlerStat struct {
	serviceURL *service.Shorten
	conf       *config.Config
}

// Service constructor.
func NewHandlerStat(conf *config.Config, short *service.Shorten) *HandlerStat {
	return &HandlerStat{serviceURL: short, conf: conf}
}

// @Summary      API get shortener statistic
// @Description  Get num of URLs and Users in the Shortener.
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      200 "OK"
// @Router       /api/internal/stats [get]
func (u *HandlerStat) GetStat(res http.ResponseWriter, req *http.Request) {
	trust := req.Context().Value(config.CtxAllow{}).(bool)
	if !trust {
		http.Error(res, "Not trusted network", http.StatusForbidden)
		return
	}

	shorts, err := u.serviceURL.GetNumShorts(req.Context())
	if err != nil {
		et := "Error during getting num of shorts (URLs): " + err.Error()
		zap.S().Errorln(et)
		http.Error(res, et, http.StatusInternalServerError)
	}
	users, err := u.serviceURL.GetNumUsers(req.Context())
	if err != nil {
		et := "Error during getting num of users: " + err.Error()
		zap.S().Errorln(et)
		http.Error(res, et, http.StatusInternalServerError)
	}
	jsonStat, err := json.Marshal(StatDTO{Shorts: shorts, Users: users})
	if err != nil {
		http.Error(res, "Error during Marshal answer stat URL", http.StatusInternalServerError)
	}

	// Set content type.
	res.Header().Add("Content-Type", "application/json")

	// Set status code 200.
	res.WriteHeader(http.StatusOK)
	res.Write(jsonStat)
}
