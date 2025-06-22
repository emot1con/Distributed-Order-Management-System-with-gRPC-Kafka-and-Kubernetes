package transport

// import (
// 	"encoding/json"
// 	"net/http"
// 	"user_service/auth"
// 	"user_service/service"

// 	"github.com/sirupsen/logrus"
// 	"golang.org/x/oauth2"
// )

// func Routes() http.Handler {
// 	r := http.NewServeMux()

// 	r.HandleFunc("/google/callback", nil)
// 	r.HandleFunc("/facebook/callback", nil)
// 	r.HandleFunc("/github/callback", nil)
// 	return r
// }

// type OauthHandler struct {
// 	service *service.UserService
// }

// func NewOauthHandler(service *service.UserService) *OauthHandler {
// 	return &OauthHandler{
// 		service: service,
// 	}
// }

// func (c *OauthHandler) OAuthFacebookCallback(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")
// 	if code == "" {
// 		http.Error(w, `{"error": "missing code"}`, http.StatusBadRequest)
// 		return
// 	}

// 	logrus.Info("usecase facebook auth")
// 	jwtToken, err := c.service.FacebookAuth(code)
// 	if err != nil {
// 		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(jwtToken)

// 	logrus.Info("facebook oauth callback successful")
// }

// func (c *OauthHandler) OAuthGithubCallback(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")
// 	if code == "" {
// 		http.Error(w, `{"error": "missing code"}`, http.StatusBadRequest)
// 		return
// 	}

// 	logrus.Info("usecase github auth")

// 	jwtToken, err := c.service.GithubAuth(code)
// 	if err != nil {
// 		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(jwtToken)

// 	logrus.Info("github oauth callback successful")
// }
