package github_openid_wrapper

import (
	"context"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/packr/v2"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
	githubauth "golang.org/x/oauth2/github"
	"strings"
	"time"
)

var GithubLoginUrl = "https://github.com/login/oauth/authorize"
var GithubAccessTokenUrl = "https://github.com/login/oauth/access_token"
var GitHubClientSecret = ""
var GitHubClientId = ""
var IDTokenIssuer = ""
var CognitoRedirectUri = "https://####.auth.eu-central-1.amazoncognito.com/oauth2/idpresponse"

var RsaPrivateKey = "jwtRS256.key"
var RsaPublicKey = RsaPrivateKey+".pub"

type AuthorizeRequest struct {
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope"`
	State        string `json:"state,omitempty"`
	ResponseType string `json:"response_type"`
}
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt int64 `json:"expires_at"`
}
type TokenRequest struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
	State string `json:"state,omitempty"`
}

func (auth *AuthorizeRequest) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, auth)
}

func (auth *AuthorizeRequest) UnmarshalMap(query map[string]string) error {
	config := &mapstructure.DecoderConfig{
		Result:   auth,
		TagName: "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil{
		return err
	}
	return decoder.Decode(query)
}

func (auth *AuthorizeRequest) GetAuthorizeUrl() string {
	scopes := strings.Split(auth.Scope, " ")
	scopes = append(scopes, oidc.ScopeOpenID)
	conf := &oauth2.Config{
		ClientID: auth.ClientID,
		Scopes: scopes,
		Endpoint: githubauth.Endpoint,
	}
	return conf.AuthCodeURL(auth.State)
}

func (tokenrequest *TokenRequest) MakeTokenRequest(ctx context.Context) ( *TokenResponse, error) {
	conf := &oauth2.Config{
		ClientID: tokenrequest.ClientID,
		ClientSecret: tokenrequest.ClientSecret,
		RedirectURL: tokenrequest.RedirectUri,
		Scopes: []string{oidc.ScopeOpenID, "user:email", "user:read"},
		Endpoint: githubauth.Endpoint,
	}
	token, err := conf.Exchange(ctx, tokenrequest.Code)
	if err != nil{
		return nil, err
	}

	resp := &TokenResponse{
		AccessToken: token.AccessToken,
		ExpiresAt: token.Expiry.Unix(),
	}

	return resp,nil
}

func loadKey(name string) ([]byte, error){
	box := packr.New("RsaBox", "../../rsaKeys")
	return box.Find(name)
}

func MakeIDToken(token TokenResponse) (string, error){
	claims := jwt.StandardClaims{
		Audience:GitHubClientId,
		Issuer: IDTokenIssuer,
		IssuedAt: time.Now().Unix(),
		ExpiresAt: token.ExpiresAt,
	}
	key, err := loadKey(RsaPrivateKey)
	if err != nil {
		return "",err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return "",err
	}
	method := jwt.SigningMethodRS256
	rawToken := jwt.NewWithClaims(method, claims)

	return rawToken.SignedString(privateKey)
}