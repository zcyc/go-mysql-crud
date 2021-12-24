package service

import (
	"database/sql"
	"encoding/json"
	"go-example/common/result"
	"log"
	"net/http"

	"go-example/dao"
	"go-example/model"

	"github.com/gorilla/mux"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {
	var (
		users []model.User
	)
	getUserListRes, err := dao.GetUserList()
	if err != nil {
		log.Println("[GetUserList] 数据库执行错误", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get user list failed!"))
		_, err := w.Write(res)
		if err != nil {
			log.Println("[GetUserList] 数据库执行错误，返回错误", err)
			return
		}
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(getUserListRes)

	for getUserListRes.Next() {
		var u model.User
		err := getUserListRes.Scan(&u.ID, &u.Name, &u.Password, &u.Status)
		if err != nil {
			log.Println("[getUserListRes.Scan] 数据库数据转对象错误", err)
		}
		users = append(users, u)
	}
	if err := getUserListRes.Err(); err != nil {
		log.Println("[getUserListRes.Err] 数据库数据转对象错误", err)
	}

	if len(users) == 0 {
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "No users"))
		_, err := w.Write(res)
		if err != nil {
			log.Println("[GetUserList] 数据库没有数据")
			return
		}
		return
	}

	res, _ := json.Marshal(result.SuccessDate(users))
	if _, err := w.Write(res); err != nil {
		log.Println("[GetUserList] 获取数据成功，返回错误", err)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		log.Println("[GetUser] 参数错误")
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get failed,id is nil!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser] 请求参数错误，返回错误", err)
			return
		}
		return
	}
	log.Println("[GetUser][id]", id)
	// 接收查询结果
	var (
		name     string
		password string
		status   int
	)

	if err := dao.GetUser(&id, &name, &password, &status); err != nil {
		log.Println("[GetUser] 数据库错误", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "User not found!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser] 数据库错误，返回错误", err)
			return
		}
		return
	}

	// 返回查询结果
	user := model.User{ID: id, Name: name, Password: password, Status: status}
	userRes := result.SuccessDate(user)
	res, err := json.Marshal(userRes)
	if err != nil {
		log.Println("[GetUser] 数据库数据转返回结果错误", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get user failed!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser] 数据库数据转返回结果错误，返回错误", err)
			return
		}
		return
	}
	if _, err := w.Write(res); err != nil {
		log.Println("[GetUser] 获取数据成功，返回错误", err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("[CreateUser] 参数转换错误", err)
		return
	}
	log.Println("[CreateUser] 参数", user)
	if user.Name == "" || user.Password == "" {
		log.Printf("[CreateUser] Name:{%s} Password:{%s} 最少有一个是空的", user.Name, user.Password)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Name/Password is empty!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[CreateUser] 参数错误，返回错误", err)
			return
		}
		return
	}
	createUserRes, err := dao.CreateUser(user)
	if err != nil {
		log.Println("[CreateUser] 数据库执行错误", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, user.Name+"already exists!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[CreateUser] 数据查询错误，返回错误", err)
			return
		}
		return
	}

	id, _ := createUserRes.LastInsertId()
	log.Println("[CreateUser] 插入的id是", id)

	// 新增用户成功的返回
	res, _ := json.Marshal(result.SuccessMsg("Create user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[CreateUser] 创建用户成功，返回错误", err)
		return
	}
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("[UpdateUser] 参数转码错误", err)
		return
	}
	log.Println("[UpdateUser] 参数", user)
	if user.ID == "" || user.Name == "" || user.Password == "" {
		log.Printf("[UpdateUser] ID:{%s} Name:{%s} Password:{%s} 最少有一个是空的", user.ID, user.Name, user.Password)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "ID/Name/Password is empty!"))
		if _, err := w.Write(res); err != nil {
			log.Printf("[UpdateUser] 参数错误，返回错误")
			return
		}
		return
	}

	updateUserRes, err := dao.UpdateUser(user)
	if err != nil {
		log.Println("[UpdateUser] 数据库执行失败", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, user.Name+" update failed"))
		if _, err := w.Write(res); err != nil {
			log.Println("[UpdateUser] 数据库执行失败，返回错误", err)
			return
		}
		return
	}

	// 修改函数是只负责修改，数据是不是存在应该由业务代码去判断。
	rowsAffected, _ := updateUserRes.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[UpdateUser] 数据库执行成功，但是没有数据被修改", rowsAffected)
	}

	// 修改成功的返回值
	res, _ := json.Marshal(result.SuccessMsg("Update user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[UpdateUser] 数据修改成功，返回错误", err)
		return
	}
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		log.Println("[DeleteUser] 参数错误")
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Delete failed,id is nil!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[DeleteUser] 参数错误，返回错误", err)
			return
		}
		return
	}
	log.Println("[DeleteUser][id]", id)
	deleteUserRes, err := dao.DeleteUser(id)
	if err != nil {
		log.Println("[DeleteUser] 数据库执行错误", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Delete failed!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[DeleteUser] 数据库执行错误，返回错误", err)
			return
		}
		return
	}
	rowsAffected, _ := deleteUserRes.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[DeleteUser] 数据库执行成功，但是没有数据被修改", rowsAffected)
	}
	res, _ := json.Marshal(result.SuccessMsg("Delete user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[DeleteUser] 数据删除成功，返回错误", err)
		return
	}
}
