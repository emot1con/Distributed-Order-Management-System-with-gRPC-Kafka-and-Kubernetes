package transport

import (
	"net/http"
	"user_service/auth"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func Routes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/oauth/google", GoogleOAuthHandler)
	r.HandleFunc("/oauth/facebook", FacebookOAuthHandler)
	r.HandleFunc("/oauth/github", GithubOAuthHandler)

	r.HandleFunc("/google/callback", nil)
	r.HandleFunc("/facebook/callback", nil)
	r.HandleFunc("/github/callback", nil)
	return r
}

func GoogleOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := auth.OauthGoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
	logrus.Info("redirecting to google oauth")
}

func FacebookOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := auth.OauthFacebookConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
	logrus.Info("redirecting to facebook oauth")
}

func GithubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := auth.OauthGithubConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
	logrus.Info("redirecting to github oauth")
}
