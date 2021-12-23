package service

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-example/dao"
	"go-example/model"
	"io"
	"log"
	"net/http"
)

// TODO: 需要一个统一返回结果的接口

func GetUserList(w http.ResponseWriter, r *http.Request) {
	var (
		rows  sql.Rows
		users []model.User
	)
	if err := dao.GetUserList(&rows); err != nil {
		log.Println("[GetUserList][dao.GetUserList]", err)
		_, err := io.WriteString(w, "Get user list failed!")
		if err != nil {
			log.Println("[GetUserList][ dao.GetUserList][io.WriteString]", err)
			return
		}
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(&rows)
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Password, &u.Status)
		if err != nil {
			log.Println(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	if len(users) == 0 {
		_, err := io.WriteString(w, "No users")
		if err != nil {
			log.Println("[GetUserList] No users")
			return
		}
		return
	}
	usersJson, err := json.Marshal(users)
	if err != nil {
		log.Println("[GetUserList][json.Marshal]", err)
		_, err := io.WriteString(w, "Get user list failed!")
		if err != nil {
			log.Println("[GetUserList][json.Marshal][io.WriteString]", err)
			return
		}
		return
	}
	if _, err := w.Write(usersJson); err != nil {
		log.Println("[GetUserList][json.Marshal][w.Write]", err)
		return
	}
	//json.NewEncoder(w) 会多一个空行，所以换用 w.Write
	//err = json.NewEncoder(w).Encode(users)
	//if err != nil {
	//	log.Println("[GetUserList][Encode(users)]", err)
	//	return
	//}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	log.Println("[GetUser][id]", id)

	// 接收查询结果
	var (
		name     string
		password string
		status   int
	)

	if err := dao.GetUser(&id, &name, &password, &status); err != nil {
		log.Println("[GetUser][dao.GetUser]", err)
		if _, err := io.WriteString(w, "User not found"); err != nil {
			log.Println("[GetUser][dao.GetUser]", err)
			return
		}
		return
	}
	user, err := json.Marshal(model.User{ID: id, Name: name, Password: password, Status: status})
	if err != nil {
		log.Println("[GetUser][json.Marshal]", err)
		_, err := io.WriteString(w, "Get user failed!")
		if err != nil {
			log.Println("[GetUser][json.Marshal][io.WriteString]", err)
			return
		}
		return
	}
	// 返回查询结果
	if _, err := w.Write(user); err != nil {
		log.Println("[GetUser][json.Marshal][w.Write]", err)
		return
	}
	//json.NewEncoder(w) 会多一个空行，所以换用 w.Write
	//if err := json.NewEncoder(w).Encode(model.User{ID: id, Name: name, Password: password, Status: status}); err != nil {
	//	log.Println("[GetUser][Encode(User)]", err)
	//	return
	//}
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
		if _, err := io.WriteString(w, "Name/Password is empty"); err != nil {
			log.Printf("[CreateUser][io.WriteString] Name:{%s} Password:{%s} 最少有一个是空的", user.Name, user.Password)
			return
		}
		return
	}
	result, err := dao.CreateUser(user)
	if err != nil {
		log.Println("[CreateUser][dao.CreateUser]", err)
		if _, err := io.WriteString(w, user.Name+" already exists"); err != nil {
			log.Println("[CreateUser][dao.CreateUser][w.Write]", err)
			return
		}
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		log.Println("[CreateUser][result.LastInsertId]", err)
	}
	// 新增用户成功的返回
	_, err = io.WriteString(w, "Create user success!")
	if err != nil {
		log.Println("[CreateUser][w.Write]", err)
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
		if _, err := io.WriteString(w, "ID/Name/Password is empty"); err != nil {
			log.Printf("[UpdateUser][io.WriteString] ID:{%s} Name:{%s} Password:{%s} 最少有一个是空的", user.ID, user.Name, user.Password)
			return
		}
		return
	}

	result, err := dao.UpdateUser(user)
	if err != nil {
		log.Println("[UpdateUser][dao.UpdateUser]", err)
		if _, err := io.WriteString(w, user.Name+" update failed"); err != nil {
			log.Println("[UpdateUser][dao.UpdateUser][w.Write]", err)
			return
		}
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[UpdateUser][rowsAffected!=1]", rowsAffected)
	}

	// 修改成功的返回值
	if _, err := io.WriteString(w, "Update user success!"); err != nil {
		log.Println("[UpdateUser][w.Write]", err)
		return
	}
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 获取路由参数
	vars := mux.Vars(r)
	var id = vars["id"]
	log.Println("[DeleteUser][id]", id)
	result, err := dao.DeleteUser(id)
	if err != nil {
		log.Println("[DeleteUser][dao.DeleteUser]", err)
		if _, err := io.WriteString(w, "delete failed"+id); err != nil {
			log.Println("[DeleteUser][dao.DeleteUser][w.Write]", err)
			return
		}
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		log.Println("[DeleteUser][rowsAffected!=1]", rowsAffected)
	}
	if _, err := io.WriteString(w, "Delete user success!"); err != nil {
		log.Println("[DeleteUser][w.Write]", err)
		return
	}
}
