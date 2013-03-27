package user

import (
	"guest/model/usermodel"
	"net/http"
	"html/template"
	"strings"
	"regexp"
	"crypto/md5"
	"io"
	"fmt"
	"github.com/astaxie/session"
      _ "github.com/astaxie/session/providers/memory"
	"log"
)

type role struct {
	NonLogin bool
	IsAdmin  bool
}

type user struct {
	userID uint32
	username string
	passwd string
	permiss uint32
}

var globalSessions, _ = session.NewManager("memory", "gosessionid", 0)

func Reg(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/nav.html", "view/reg.html")

	if r.Method == "POST" {
		user := user{}
		tmp_username := strings.Trim(r.FormValue("username"), " ") 
		tmp_username = template.HTMLEscapeString(tmp_username)
		if len(tmp_username) < 5 || len(tmp_username) > 50 {
			mainData(w, r, t, "reg", "用户名少于5位或大于50位")
			return
		}
		if m,_ := regexp.MatchString("^[a-zA-z]+[a-zA-Z0-9]*$", tmp_username); !m {
			mainData(w, r, t, "reg", "用户名只能以字母和数字的组合")
			return
		}
		user.username = tmp_username

		user.passwd = r.FormValue("passwd")
		if len(user.passwd) < 5 || len(user.passwd) > 50 {
			mainData(w, r, t, "reg", "用户密码少于5位或大于50位")
			return
		}
		passwd2 := r.FormValue("passwd2")
		if user.passwd != passwd2 {
			mainData(w, r, t, "reg", "用户密码要一致")
			return
		}
		if usermodel.IsDuplicateUser(user.username) {
			mainData(w, r, t, "reg", "用户已存在")
			return
		}
		
		h := md5.New()
		salt := "olala!@#$"
		io.WriteString(h, user.passwd+salt)
		user.passwd = fmt.Sprintf("%x", h.Sum(nil))
		user.permiss = 2
		if usermodel.AddGuestUser(user.username, user.passwd, user.permiss) > 0 {
			mainData(w, r, t, "reg", "添加用户成功")
			return
		}
		mainData(w, r, t, "reg", "无法添加用户")
		return 
	} 
	mainData(w, r, t, "reg", nil)
}
func Login(w http.ResponseWriter, r *http.Request) { 
	t, _ := template.ParseFiles("view/nav.html", "view/login.html")
	
	if r.Method == "POST" {
		user := user{}	

		user.username = strings.Trim(r.FormValue("username"), " ")
		user.username = template.HTMLEscapeString(user.username)
		if len(user.username) < 5 || len(user.username) > 50 {
			mainData(w, r, t, "login", "用户名少于5位或大于50位")
			return
		}
		if m,_ := regexp.MatchString("^[a-zA-z]+[a-zA-Z0-9]*$", user.username); !m {
			mainData(w, r, t, "login", "用户名只能以字母和数字的组合")
			return
		}

		user.passwd = r.FormValue("passwd")
		if len(user.passwd) < 5 || len(user.passwd) > 50 {
			mainData(w, r, t, "login", "用户密码少于5位或大于50位")
			return
		}
		h := md5.New()
		salt := "olala!@#$"
		io.WriteString(h, user.passwd+salt)
		user.passwd = fmt.Sprintf("%x", h.Sum(nil))
	
		rows := usermodel.CheckLoginUser(user.username, user.passwd) 
		if rows.Next() {
			err := rows.Scan(&user.userID, &user.username, &user.permiss)
			checkErr(err)
			
			sess := globalSessions.SessionStart(w, r)
			sess.Set("userID", user.userID)
			sess.Set("username", user.username)
			sess.Set("userpermiss", user.permiss)
			http.Redirect(w, r, "http://localhost:9090/index", 302)
		} else {
			mainData(w, r, t, "login", "用户不存在或者密码错误")
			return
		}
	} else {
		mainData(w, r, t, "login", nil)
		return
	}
}


func Exit(w http.ResponseWriter, r *http.Request) { 
	globalSessions.SessionDestroy(w, r)
	http.Redirect(w, r, "http://localhost:9090/", 302)
}

func roleManage(w http.ResponseWriter, r *http.Request) *role {
	role := role{}

	sess := globalSessions.SessionStart(w, r)
	permiss := sess.Get("userpermiss")
	if permiss == nil {
		role.NonLogin = true
	} else if permiss.(uint32) == 1 {
		role.IsAdmin = true
	}
	return &role
}
func mainData(w http.ResponseWriter, r *http.Request, t *template.Template, template string, message interface{}) {
	role := roleManage(w, r) 
	t.ExecuteTemplate(w, "nav", role)
	t.ExecuteTemplate(w, template, message)
}


func checkErr(err error) {
    if err != nil {
	log.Fatal(err)
    }
}


