package auth

import (
	"time"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	// EXPIRATION_DELTA token timeout
	EXPIRATION_DELTA   = 3009999999 // Unit: second
	// NOT_BEFORE_DELTA ..
	NOT_BEFORE_DELTA   = 60         // Unit: second
	// NOT_BEFORE_DELTA secret key
	APIAUTH_SECRET_KEY = "bluelotus"
)

// AnsibleJwtClaims jwt struct
type AnsibleJwtClaims struct {
	UID string `json:"identity,omitempty"`
	jwt.StandardClaims
	Meta interface{} `json:"meta,omitempty"`
}

// GetUid get uuid
func (claims *AnsibleJwtClaims) GetUid() string {
	return claims.UID
}

// MetaData get metadata
func (claims *AnsibleJwtClaims) MetaData() interface{} {
	return claims.Meta
}

// IssueTokenUsingDgrijalva auth token
func IssueTokenUsingDgrijalva(uid string, meta interface{}) (string, error) {
	iat := time.Now()
	deferIat := iat.Add(-NOT_BEFORE_DELTA * time.Second)
	exp := iat.Add(EXPIRATION_DELTA * time.Second)
	nbf := iat.Add(-NOT_BEFORE_DELTA * time.Second)

	claims := AnsibleJwtClaims{
		uid,
		jwt.StandardClaims{
			IssuedAt:  deferIat.Unix(),
			ExpiresAt: exp.Unix(),
			NotBefore: nbf.Unix(),
		},
		meta,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := t.SignedString([]byte(APIAUTH_SECRET_KEY))

	return tokenStr, err

}

// JwtAuthFilter filter 
func JwtAuthFilter(ctx *context.Context) {
	if ctx.Request.RequestURI == "/ansible/login" {
		return
	}
	token, err := request.ParseFromRequestWithClaims(ctx.Request,
		request.AuthorizationHeaderExtractor,
		&AnsibleJwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(APIAUTH_SECRET_KEY), nil
		})
	if err != nil || token.Claims.(*AnsibleJwtClaims).UID == "" {
		ctx.Output.Status = 401
		ctx.Output.JSON(err, false, false)
		return
	}

	uid := token.Claims.(*AnsibleJwtClaims).UID
	ctx.Input.SetData("uid", uid)
}

