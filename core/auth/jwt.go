package auth

import (
	"time"
	"fmt"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

)

const (
	EXPIRATION_DELTA   = 3009999999 // Unit: second
	NOT_BEFORE_DELTA   = 60  // Unit: second
	APIAUTH_SECRET_KEY = "bluelotus"
)

type AnsibleJwtClaims struct {
	Uid string `json:"identity,omitempty"`
	jwt.StandardClaims
	Meta interface{} `json:"meta,omitempty"`
}

func (claims *AnsibleJwtClaims) GetUid() string {
	return claims.Uid
}

func (claims *AnsibleJwtClaims) MetaData() interface{} {
	return claims.Meta
}

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


func JwtAuthFilter(ctx *context.Context) {
	token, err := request.ParseFromRequestWithClaims(ctx.Request,
		request.AuthorizationHeaderExtractor,
		&AnsibleJwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(APIAUTH_SECRET_KEY), nil
		})
	if err != nil || token.Claims.(*AnsibleJwtClaims).Uid == "" {
		fmt.Println(err)
		ctx.Output.Status = 401
		ctx.Output.JSON(err, false, false)
		return
	}
	
		uid := token.Claims.(*AnsibleJwtClaims).Uid
		ctx.Input.SetData("uid", uid)
	}