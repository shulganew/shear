package rest

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
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
	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectURL, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}
	zap.S().Infoln("redirectURL ", redirectURL)
	brief := service.GenerateShortLinkByte()
	mainURL, answerURL, err := u.serviceURL.GetAnsURLFast(redirectURL.Scheme, u.conf.GetResponse(), brief)
	if err != nil {
		http.Error(res, "Error parse URL", http.StatusInternalServerError)
		return
	}

	// Set content type.
	res.Header().Add("Content-Type", "text/plain")

	// find UserID in cookies
	userID, err := req.Cookie("user_id")
	if err != nil {
		http.Error(res, "Can't find user in cookies", http.StatusUnauthorized)
	}

	// Save map to storage.
	err = u.serviceURL.AddURL(req.Context(), userID.Value, brief, (*redirectURL).String())

	if err != nil {
		var tagErr *service.ErrDuplicatedURL
		if errors.As(err, &tagErr) {
			// set status code 409 Conflict
			res.WriteHeader(http.StatusConflict)

			// send existed string from error
			var answer string
			answer, err = url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}
			res.Write([]byte(answer))
			return
		}

		zap.S().Errorln(err)
		http.Error(res, "Error saving in Storage.", http.StatusInternalServerError)
		return
	}
	// Set status code 201.
	res.WriteHeader(http.StatusCreated)
	// Send generate and saved string.
	res.Write([]byte(answerURL.String()))
}
