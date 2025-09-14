package types

import "encoding/json"

type User struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	FullName  string `gorm:"type:varchar(100);not null" json:"full_name" validate:"required,min=3,max=100"`
	Email     string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string `gorm:"not null" json:"-" validate:"required,min=6,max=100"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RegisterPayload struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`
}

type TokenResponse struct {
	Message               string `json:"message"`
	Token                 string `json:"token"`
	ExpiredAt             string `json:"expired_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at"`
	Role                  string `json:"role"`
}

type OAuthUserData struct {
	ProviderID string `json:"provider_id"`
	Provider   string `json:"provider"` // "google", "github", "facebook"
	Email      string `json:"email"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

type OauthGithubUserModel struct {
	ProviderID int64  `json:"id"`
	Provider   string `json:"provider"`
	Login      string `json:"login"`
	AvatarURL  string `json:"avatar_url"`
	Email      string `json:"email"`
}

// OAuthState stores state information for OAuth flow
type OAuthState struct {
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
	ExpiresAt   int64  `json:"expires_at"`
}

// OAuthConfig contains OAuth provider configuration
type OAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	UserInfoURL  string `json:"user_info_url"`
	Scopes       string `json:"scopes"`
}

func (u *OAuthUserData) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Ambil provider ID dari "sub" atau "id"
	if sub, ok := raw["sub"].(string); ok {
		u.ProviderID = sub
	} else if id, ok := raw["id"].(string); ok {
		u.ProviderID = id
	}

	if provider, ok := raw["provider"].(string); ok {
		u.Provider = provider
	}
	if email, ok := raw["email"].(string); ok {
		u.Email = email
	}
	if name, ok := raw["name"].(string); ok {
		u.Name = name
	}
	if avatar, ok := raw["picture"].(string); ok {
		u.AvatarURL = avatar
	} else if avatar, ok := raw["avatar_url"].(string); ok {
		u.AvatarURL = avatar
	}

	return nil
}
