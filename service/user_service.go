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
		log.Println("[GetUserList][dao.GetUserList]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get user list failed!"))
		_, err := w.Write(res)
		if err != nil {
			log.Println("[GetUserList][ dao.GetUserList][json.Marshal]", err)
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
			log.Println(err)
		}
		users = append(users, u)
	}
	if err := getUserListRes.Err(); err != nil {
		log.Println(err)
	}

	if len(users) == 0 {
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "No users"))
		_, err := w.Write(res)
		if err != nil {
			log.Println("[GetUserList] No users")
			return
		}
		return
	}

	res, _ := json.Marshal(result.SuccessDate(users))
	if _, err := w.Write(res); err != nil {
		log.Println("[GetUserList][json.Marshal][w.Write]", err)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		log.Println("[GetUser][id] nil")
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get failed,id is nil!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser][dao.GetUser][w.Write]", err)
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
		log.Println("[GetUser][dao.GetUser]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "User not found!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser][dao.GetUser]", err)
			return
		}
		return
	}

	// 返回查询结果
	user := model.User{ID: id, Name: name, Password: password, Status: status}
	userRes := result.SuccessDate(user)
	res, err := json.Marshal(userRes)
	if err != nil {
		log.Println("[GetUser][json.Marshal]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Get user failed!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[GetUser][dao.GetUser]", err)
			return
		}
		return
	}
	if _, err := w.Write(res); err != nil {
		log.Println("[GetUser][json.Marshal][w.Write]", err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("[CreateUser][Decode(&user)]", err)
		return
	}
	log.Println("[CreateUser][Request]", user)
	if user.Name == "" || user.Password == "" {
		log.Printf("[CreateUser] Name:{%s} Password:{%s} 最少有一个是空的", user.Name, user.Password)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Name/Password is empty!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[CreateUser][json.Marshal]", err)
			return
		}
		return
	}
	createUserRes, err := dao.CreateUser(user)
	if err != nil {
		log.Println("[CreateUser][dao.CreateUser]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, user.Name+"already exists!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[CreateUser][dao.CreateUser]", err)
			return
		}
		return
	}

	_, err = createUserRes.LastInsertId()
	if err != nil {
		log.Println("[CreateUser][result.LastInsertId]", err)
	}
	// 新增用户成功的返回
	res, _ := json.Marshal(result.SuccessMsg("Create user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[CreateUser][json.Marshal]", err)
		return
	}
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("[CreateUser][Decode(&user)]", err)
		return
	}
	log.Println("[UpdateUser][Request]", user)
	if user.ID == "" || user.Name == "" || user.Password == "" {
		log.Printf("[UpdateUser] ID:{%s} Name:{%s} Password:{%s} 最少有一个是空的", user.ID, user.Name, user.Password)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "ID/Name/Password is empty!"))
		if _, err := w.Write(res); err != nil {
			log.Printf("[UpdateUser][json.Marshal][nil][w.Write]异常")
			return
		}
		return
	}

	updateUserRes, err := dao.UpdateUser(user)
	if err != nil {
		log.Println("[UpdateUser][dao.UpdateUser]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, user.Name+" update failed"))
		if _, err := w.Write(res); err != nil {
			log.Println("[UpdateUser][dao.UpdateUser][w.Write]", err)
			return
		}
		return
	}

	rowsAffected, _ := updateUserRes.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[UpdateUser][rowsAffected!=1]", rowsAffected)
	}

	// 修改成功的返回值
	res, _ := json.Marshal(result.SuccessMsg("Update user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[UpdateUser][w.Write]", err)
		return
	}
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		log.Println("[DeleteUser][id] nil")
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Delete failed,id is nil!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[DeleteUser][dao.DeleteUser][w.Write]", err)
			return
		}
		return
	}
	log.Println("[DeleteUser][id]", id)
	deleteUserRes, err := dao.DeleteUser(id)
	if err != nil {
		log.Println("[DeleteUser][dao.DeleteUser]", err)
		res, _ := json.Marshal(result.FailedMsg(result.ERROR_USER, "Delete failed!"))
		if _, err := w.Write(res); err != nil {
			log.Println("[DeleteUser][dao.DeleteUser][w.Write]", err)
			return
		}
		return
	}
	rowsAffected, _ := deleteUserRes.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[DeleteUser][rowsAffected!=1]", rowsAffected)
	}
	res, _ := json.Marshal(result.SuccessMsg("Delete user success!"))
	if _, err := w.Write(res); err != nil {
		log.Println("[DeleteUser][w.Write]", err)
		return
	}
}
