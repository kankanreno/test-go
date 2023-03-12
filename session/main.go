package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

// 存储所有用户的session数据
var SessionMap map[string]*sessions.Session

func set(w http.ResponseWriter, r *http.Request) {
	// get sessionId
	var sessionId string
	if cookie, err := r.Cookie("session_id"); err != nil {
		sessionId = base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	} else {
		sessionId = cookie.Value
	}

	// gen session
	var session *sessions.Session
	var ok bool
	if session, ok = SessionMap[sessionId]; !ok {
		session = &sessions.Session{
			ID:     sessionId,
			Values: make(map[any]any),
		}
	}
	session.Values["account"] = "kankan"
	session.Values["name"] = "看看"

	// session add to sessionMap
	SessionMap[sessionId] = session

	// sessionId set to cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Path:    "/",
		Domain:  "",
		Expires: time.Now().Add(time.Hour * 24),
	})

	fmt.Fprintln(w, "Hello World")
}

func read(w http.ResponseWriter, r *http.Request) {
	// get sessionId
	var sessionId string
	if cookie, err := r.Cookie("session_id"); err != nil {
		jsonErrorHandler(w, 1, err.Error(), nil)
		return
	} else {
		sessionId = cookie.Value
	}

	// get and return session
	var session *sessions.Session
	var ok bool
	if session, ok = SessionMap[sessionId]; !ok {
		jsonErrorHandler(w, 1, "session 不存在", nil)
		return
	}
	fmt.Fprintf(w, "account: %s, name: %s", session.Values["account"], session.Values["name"])
}

func main() {
	SessionMap = make(map[string]*sessions.Session)

	m := http.NewServeMux()
	m.HandleFunc("/set", set)
	m.HandleFunc("/read", read)

	server := &http.Server{
		Addr:    ":9999",
		Handler: m,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error from HTTP Server: %v", err)
	}
}

// 返回JSON结果
func jsonErrorHandler(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
