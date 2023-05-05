package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/session"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserInfo() map[string]interface{} {
	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	email := strings.TrimSpace(strings.Split(user[3], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[4], "=")[1])

	return map[string]interface{}{
		"username":    username,
		"displayName": displayName,
		"email":       email,
		"photoURL":    photoURL,
	}
}

func Inituser(userVO *vo.UserVO) {
	username := "username=" + userVO.Username
	password := "password=" + userVO.Password
	displayName := ""
	email := ""
	photoURL := ""
	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{username, password, displayName, photoURL, email})
}

func Login(userVO *vo.UserVO, ctx *gin.Context, sessions *session.Manager) *vo.Response {

	response := &vo.Response{}

	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	// username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	email := strings.TrimSpace(strings.Split(user[3], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[4], "=")[1])

	if username != userVO.Username || password != userVO.Password {
		log.Panicln("User authentication failed")
		response.Code = 401
		response.Msg = "User authentication failed"
		return response
	}

	session := sessions.Start(ctx.Writer, ctx.Request)

	session.Set("username", username)

	userVO.SessionID = session.SessionID()
	response.Code = 200
	response.Msg = "Login success"
	userVO.Password = ""
	response.Data = map[string]interface{}{
		"username":    username,
		"displayName": displayName,
		"email":       email,
		"photoURL":    photoURL,
	}

	return response
}

func Logout(ctx *gin.Context, sessions *session.Manager) {
	sessions.Destroy(ctx.Writer, ctx.Request)
}

func ChangeUser(username, password string) {
	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	email := strings.TrimSpace(strings.Split(user[3], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[4], "=")[1])
	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{
		"username = " + username,
		"password = " + password,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
		"email=" + email,
	})
}

func ChangePassword(newPassword string) *vo.Response {

	response := &vo.Response{}
	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	email := strings.TrimSpace(strings.Split(user[3], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[4], "=")[1])
	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{
		"username = " + username,
		"password = " + newPassword,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
		"email=" + email,
	})

	response.Code = 200
	response.Msg = "Update user new password success"

	return response
}
