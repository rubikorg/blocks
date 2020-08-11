package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	r "github.com/rubikorg/rubik"
	"github.com/rubikorg/rubik/pkg"
)

// BlockName is the name of this rubik block
const BlockName = "JWT"

type config struct {
	Secret              string `json:"secret"`
	CookieHTTPOnly      bool   `json:"cooke_http_only"`
	CookieHTTPOnlyError string `json:"cooke_http_only_error"`
	CookieKey           string `json:"cookie_key"`
	UnauthorizedError   string `json:"unauth_error"`
	Expiry              int    `json:"expiry_time"`
	ExpiryClaimsKey     string `json:"expiry_key"`
}

var blockConfig config

// BlockJWT is the block for jwt authentication
type BlockJWT struct{}

// CreateToken creates a jwt token for given secret
func CreateToken(claims jwt.MapClaims, expiry bool) (string, error) {
	var err error
	if expiry && blockConfig.Expiry != 0 {
		key := blockConfig.ExpiryClaimsKey
		if key == "" {
			key = "exp"
		}
		claims[key] = time.Now().Add(time.Duration(blockConfig.Expiry) * time.Second).Unix()
	} else if expiry && blockConfig.Expiry == 0 {
		pkg.ErrorMsg(
			"need to set `expiry_time` inside config to set expiry(in seconds) inside token")
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(blockConfig.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// CookieJWTMiddleware authorizes your request based on the cookie `token`
// presence
func CookieJWTMiddleware(req *r.Request) {
	key := blockConfig.CookieKey
	if key == "" {
		key = "token"
	}
	cookie, err := req.Raw.Cookie(key)
	if err != nil {
		req.Throw(http.StatusUnauthorized, err, r.Type.JSON)
		return
	}
	if blockConfig.CookieHTTPOnly && !cookie.HttpOnly {
		msg := blockConfig.CookieHTTPOnlyError
		if msg == "" {
			msg = "Broken cookie through HTTP"
		}
		req.Throw(http.StatusUnauthorized, r.E(msg), r.Type.JSON)
		return
	}

	claims, statusCode, err := verifyClaims(cookie.Value)
	if err != nil {
		req.Throw(statusCode, err)
		return
	}

	req.Claims = claims
}

// HeaderJWTMiddleware authorizes your request based on the cookie `token`
// presence
func HeaderJWTMiddleware(req *r.Request) {
	token := getJWTFromHeader(req.Raw)
	if token == "" {
		req.Throw(http.StatusForbidden, r.E("Token not found"), r.Type.JSON)
		return
	}
	claims, statusCode, err := verifyClaims(token)
	if err != nil {
		req.Throw(statusCode, err, r.Type.JSON)
		return
	}

	req.Claims = claims
}

// ParseClaimsMiddleware lets your parse the claims if the token is present in
// the Authorization header. The middleware does not throw any error on the
// absence of a token, so if the token is present then req.Claims will have
// your claims
func ParseClaimsMiddleware(req *r.Request) {
	token := getJWTFromHeader(req.Raw)
	if token != "" {
		claims, _ := getClaimsFromTokenString(token)
		req.Claims = claims
	}
}

func getJWTFromHeader(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	if splitted := strings.Split(bearToken, " "); len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func getClaimsFromTokenString(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Method.Alg())
		}

		return []byte(blockConfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, r.E("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, r.E("token not created with JWT block")
	}

	return claims, nil
}

func verifyClaims(tokenStr string) (jwt.MapClaims, int, error) {
	claims, err := getClaimsFromTokenString(tokenStr)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	// if expiry in seconds is set then check the exp
	if blockConfig.Expiry != 0 {
		key := blockConfig.ExpiryClaimsKey
		if key == "" {
			key = "exp"
		}
		timestamp, ok := claims[key].(int64)
		if ok {
			expTime := time.Unix(timestamp, 0)
			if time.Now().Sub(expTime) > 0 {
				return claims, http.StatusConflict, r.E("Token Expired")
			}
		}
	}

	return claims, 0, nil
}

// OnAttach implements the blocks interface
func (b BlockJWT) OnAttach(app *r.App) error {
	if err := app.Decode("jwt_auth", &blockConfig); err != nil {
		return err
	}
	return nil
}

func init() {
	r.Attach(BlockName, BlockJWT{})
}
