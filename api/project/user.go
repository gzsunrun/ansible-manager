package project

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/config"
)

type MyCustomClaims struct {
	UserID int
	jwt.StandardClaims
}

//login
func Login(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	password := r.FormValue("password")
	if account == "" || password == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h := md5.New()
	io.WriteString(h, password)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	var user db.User
	res, err := db.MysqlDB.Table("ansible_user").Where("user_account=? and user_password=?", account, passwdMD5).Get(&user)
	if !res || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := MyCustomClaims{
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "ansible-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(config.Cfg.AnsibleManager.JwtSecret))
	msg := map[string]interface{}{
		"token":   signedToken,
		"account": account,
	}
	JsonWrite(w, 200, msg)

}

//check token
func Auth(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("Auth")
	if err != nil {
		return errors.New("err")
	}

	splitCookie := strings.Split(cookie.String(), "Auth=")

	token, err := jwt.ParseWithClaims(splitCookie[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(config.Cfg.AnsibleManager.JwtSecret), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		context.Set(r, "Claims", claims)
	} else {
		return errors.New("err")
	}
	return nil
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	var user db.User
	_, err := db.MysqlDB.Table("ansible_user").Where("user_id", claims.UserID).Get(&user)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	JsonWrite(w, 200, user)
	context.Clear(r)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("user_account")
	password := r.FormValue("user_password")
	if account == "" || password == "" {
		w.WriteHeader(404)
		return
	}
	h := md5.New()
	io.WriteString(h, password)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	_, err := db.MysqlDB.Exec("insert into ansible_user (user_account,user_password) values (?,?)", account, passwdMD5)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var users []db.User
	err := db.MysqlDB.Table("ansible_user").Find(&users)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	JsonWrite(w, 200, users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	var user db.User
	_, err := db.MysqlDB.Table("ansible_user").Where("user_id=?", userID).Get(&user)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	JsonWrite(w, 200, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("user_account")
	password := r.FormValue("user_password")
	id := r.FormValue("user_id")
	if account == "" || password == "" || id == "" {
		w.WriteHeader(404)
		return
	}
	h := md5.New()
	io.WriteString(h, password)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	_, err := db.MysqlDB.Exec("update ansible_user set user_account=?,user_password=? where user_id=?", account, passwdMD5, id)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
