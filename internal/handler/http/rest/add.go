package rest

import (
	"io"
	"net/http"

	"github.com/shulganew/shear.git/internal/builders"
	"github.com/shulganew/shear.git/internal/config"
)

// POTS and add generate short URL.
// @Summary      Set origin URL
// @Description  set URL in body POST
// @Tags         api
// @Accept       plain
// @Produce      plain
// @Success      201 {string}  string  "Created"
// @Failure      401 "User unauthorized"
// @Failure      404 "Conflict. URL existed."
// @Failure      500 "Handling error"
// @Router       / [post]
func (u *HandlerURL) AddURL(res http.ResponseWriter, req *http.Request) {
	// get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)

	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	// Get resDTO from service.
	resDTO := u.serviceURL.AddURL(req.Context(), builders.AddRequestDTO{Origin: string(readBody), CtxConfig: ctxConfig, Resp: u.conf.GetResponse()})

	// Set content type.
	res.Header().Add("Content-Type", "text/plain")
	// Set status code
	res.WriteHeader(resDTO.Status.GetStatusREST())
	// Send generate and saved string.
	res.Write([]byte(resDTO.AnwerURL))
}
