package content

import (
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
	"github.com/dchest/captcha"
	"guest/model/contentmodel"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type content struct {
	Id            uint32
	Content       string
	UserID        uint32
	Username      string
	Time          string
	WhetherDelete bool
}
type page struct {
	perpage     int
	CurrentPage int
	TotalPage   int
	PrevPage    int
	NextPage    int
}
type role struct {
	NonLogin bool
   	IsAdmin  bool
}

type index struct {
	Content []*content
	page
	Message string
	CaptchaId string
}

var globalSessions, _ = session.NewManager("memory", "gosessionid", 0)

func GetGuest(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/nav.html", "view/index.html")
	index := allGuest(w, r)
	mainData(w, r, t, "index", index)
	return
}

func AddGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sess := globalSessions.SessionStart(w, r)
		if sess.Get("userID") == nil {
			http.Redirect(w, r, "http://localhost:9090/login", 302)
		}

		t, _ := template.ParseFiles("view/index.html", "view/nav.html")

		index := allGuest(w, r)

		content := content{}

		guest := strings.Trim(r.FormValue("content"), " ")
		guest = template.HTMLEscapeString(guest)
		if len(guest) == 0 || len(guest) > 500 {
			index.Message = "留言不能为空或者超过500"
			mainData(w, r, t, "index", index)
			return
		}
		captchavalue := r.FormValue("captchaSolution")
		captchaId := r.FormValue("captchaId")
		if len(captchavalue) == 0 || len(captchavalue) > 6 {
			index.Message = "验证码不能为空或者超过6位"
			mainData(w, r, t, "index", index)
			return
		}
		if !captcha.VerifyString(captchaId, captchavalue) {
			index.Message = "验证码输入错误"
			mainData(w, r, t, "index", index)
			return
		}
		content.Content = template.HTMLEscapeString(guest)
		userID := sess.Get("userID")
		content.UserID = userID.(uint32)
		if contentmodel.AddGuest(content.Content, content.UserID) > 0 {
			http.Redirect(w, r, "http://localhost:9090/", 302)
		}
		
		index.Message = "无法留言"
		mainData(w, r, t, "index", index)
		return
	}
	http.Redirect(w, r, "http://localhost:9090/login", 302)
}

func DelGuest(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if sess.Get("userID") == nil {
		http.Redirect(w, r, "http://localhost:9090/login", 302)
	}
	
	index := allGuest(w, r)
	t, _  := template.ParseFiles("view/index.html", "view/nav.html")	

	userid := sess.Get("userID")
	userID := userid.(uint32)
	id := r.FormValue("id")
	id = template.HTMLEscapeString(id)	
	contentid, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "http://localhost:9090", 302)
		return
	}
	if ! contentmodel.CheckUserID(contentid, userID).Next() {
		index.Message = "请不要乱来"
		mainData(w, r, t, "index", index)
		return
	}
	if contentmodel.DelGuest(contentid) > 0 {
		http.Redirect(w, r, "http://localhost:9090", 302)
	} else {
		index.Message = "删除留言失败"
		t.ExecuteTemplate(w, "index", index)
		return
	}
}

func doPage(r *http.Request) *index {
	index := index{}

	index.perpage = 3
	totalcount := 0
	rows := contentmodel.GetTotalCount()
	if rows.Next() {
		rows.Scan(&totalcount)
	}
	if totalcount%index.perpage > 0 {
		index.TotalPage = (totalcount / index.perpage) + 1
	} else {
		index.TotalPage = totalcount / index.perpage
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if page <= 0 || err != nil {
		index.CurrentPage = 1
	} else {
		index.CurrentPage = page
	}

	if index.CurrentPage > index.TotalPage {
		index.CurrentPage = index.TotalPage
	}
	index.PrevPage = index.CurrentPage - 1
	if index.PrevPage <= 0 {
		index.PrevPage = 1
	}
	index.NextPage = index.CurrentPage + 1
	if index.NextPage > index.TotalPage {
		index.NextPage = index.CurrentPage
	}
	return &index
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

func allGuest(w http.ResponseWriter, r *http.Request) *index {	
	index := doPage(r)

	sess := globalSessions.SessionStart(w, r)
	userID := sess.Get("userID")	
	tmp_content := content{}
	offset := (index.CurrentPage - 1) * index.perpage
	rows := contentmodel.GetGuest(offset, index.perpage)
	for rows.Next() {
		rows.Scan(&tmp_content.Id, &tmp_content.Content, &tmp_content.UserID, &tmp_content.Username, &tmp_content.Time)
		tmp_content.WhetherDelete = false
		if (userID != nil) && (userID == tmp_content.UserID) {
			tmp_content.WhetherDelete = true
		}
		index.Content = append(index.Content, &tmp_content)
	}
	index.CaptchaId = captcha.New()
	return index
}
