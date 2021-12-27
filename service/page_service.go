package service

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go-example/common/session"
	"go-example/dao"
	"go-example/model"
)

func Index(w http.ResponseWriter, r *http.Request) {
	user, _ := session.GetSession(w, r).GetAttr("user")

	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		log.Println("[index] 加载主页模版错误", err)
	}

	err = t.Execute(w, user)
	if err != nil {
		log.Println("[index] 生成主页错误", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		loginHTML, err := ioutil.ReadFile("html/login.html")
		if err != nil {
			log.Println("[login] 加载登陆页模版错误", err)
		}

		_, err = w.Write(loginHTML)
		if err != nil {
			log.Println("[login] 生成登陆页错误", err)
		}

		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	log.Println("login", name, password)
	if isEmpty(name, password) {
		Message(w, r, "字段不能为空")
		return
	}

	user, err := dao.GetUserByNameAndPassword(name, password)
	if err != nil {
		log.Println("[login] 数据库执行错误", err)
		Message(w, r, "登录失败！")
		return
	}
	if !user.Next() {
		Message(w, r, "登录失败！")
		return
	}
	log.Println("[Login] 登陆成功，即将跳转")

	// 登陆成功
	var u model.User
	err = user.Scan(&u.ID, &u.Name, &u.Password, &u.Status)
	if err != nil {
		log.Println("[Login] 登陆成功，提取页面生成参数报错")
		return
	}
	sess := session.GetSession(w, r)
	sess.SetAttr("user", u)
	http.Redirect(w, r, "/", 302)
}

func Message(w http.ResponseWriter, r *http.Request, message string) {
	t, err := template.ParseFiles("html/message.html")
	if err != nil {
		log.Println("[Message] 加载消息模版失败", err)
	}

	err = t.Execute(w, map[string]string{"Message": message})
	if err != nil {
		log.Println("[Message] 生成消息页失败", err)
	}
}

func Userinfo(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)
	user, exist := sess.GetAttr("user")
	if !exist {
		http.Redirect(w, r, "/", 302)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("html/userinfo.html")
		if err != nil {
			log.Println("[UserInfo] 加载用户也模版失败", err)
		}
		err = t.Execute(w, user)
		if err != nil {
			log.Println("[UserInfo] 生成用户页失败", err)
		}
		return
	}

	// POST 更新用户信息
	id := r.FormValue("id")
	name := r.FormValue("name")
	password := r.FormValue("password")
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		log.Println("[UserInfo] status 转换错误", err)
		return
	}
	log.Printf("[Userinfo][Update] id:%s,name:%s,password:%s,status:%d", id, name, password, status)
	if isEmpty(id, name, password) {
		Message(w, r, "字段不能为空")
		return
	}

	userUpdate := model.User{ID: id, Name: name, Password: password, Status: status}
	_, err = dao.UpdateUser(userUpdate)
	if err != nil {
		return
	}
	sess.SetAttr("user", userUpdate)
	http.Redirect(w, r, "/userinfo", 302)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		registerHTML, err := ioutil.ReadFile("html/register.html")
		if err != nil {
			log.Println("[Register] 加载注册模版失败", err)
		}
		_, err = w.Write(registerHTML)
		if err != nil {
			log.Println("[Register] 生成注册页面失败", err)
			return
		}
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		log.Println("[UserInfo] status 转换错误", err)
		return
	}

	if isEmpty(name, password, password2) {
		Message(w, r, "字段不能为空")
		return
	}

	if password != password2 {
		Message(w, r, "两次密码不相符")
		return
	}

	user := model.User{
		Name:     name,
		Password: password,
		Status:   status,
	}
	_, err = dao.CreateUser(user)
	if err != nil {
		return
	}
	Message(w, r, "注册成功！")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)
	sess.DelAttr("user")
	http.Redirect(w, r, "/", 302)
}

func isEmpty(strSlice ...string) (isEmpty bool) {
	for _, str := range strSlice {
		str = strings.TrimSpace(str)
		if str == "" || len(str) == 0 {
			isEmpty = true
			return
		}
	}
	isEmpty = false
	return
}
