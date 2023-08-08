package controller

import (
	"encoding/json"
	"flexpoint/connection"
	"flexpoint/helper"
	"flexpoint/model"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ResponseJson = helper.ResponseJson
var ResponseError = helper.ResponseError

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var user model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	// check if email already exists
	existingUser := model.User{}
	if err := connection.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		ResponseError(w, http.StatusBadRequest, "Email already exists")
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = string(hashedPassword)

	query := `INSERT INTO users (role, name, email, password)
			  VALUES (?, ?, ?, ?)`

	if err := connection.DB.Exec(query, user.Role, user.Name, user.Email, user.Password).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := connection.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	queryPoint := `INSERT INTO points (user_id, amount)
				   VALUES (?, ?)`
	if err := connection.DB.Exec(queryPoint, user.Id, 0).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := make(map[string]interface{})
	responseData["message"] = "Success Create User"

	helper.ResponseJson(w, http.StatusCreated, responseData)

}

func GetListUser(w http.ResponseWriter, r *http.Request) {
	var users []model.User

	query := r.URL.Query()
	role := query.Get("role")
	pageStr := query.Get("page")
	pageSizeStr := query.Get("page_size")

	// Default pagination values
	page := 1
	pageSize := 10

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			ResponseError(w, http.StatusBadRequest, "Invalid value for 'page' parameter")
			return
		}
	}

	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			ResponseError(w, http.StatusBadRequest, "Invalid value for 'page_size' parameter")
			return
		}
	}

	offset := (page - 1) * pageSize

	db := connection.DB
	if role != "" {
		if role == "admin" {
			db = db.Where("role_id = ?", 1001)
		}
		if role == "user" {
			db = db.Where("role_id = ?", 1002)
		}
	}

	var totalCount int64
	if err := db.Model(&model.User{}).Count(&totalCount).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var totalPages int64
	if totalCount > 0 {
		totalPages = int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	} else {
		totalPages = 0
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := make(map[string]interface{})
	meta["page"] = page
	meta["page_size"] = pageSize
	meta["total"] = totalCount
	meta["total_pages"] = totalPages

	responseData := make(map[string]interface{})
	var data []interface{}
	for _, user := range users {
		item := struct {
			ID    int    `json:"id"`
			Role  string `json:"role"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			ID:    user.Id,
			Role:  user.Role,
			Name:  user.Name,
			Email: user.Email,
		}
		data = append(data, item)
	}

	responseData["data"] = data
	responseData["meta"] = meta
	responseData["message"] = "Success Get All Users"

	helper.ResponseJson(w, http.StatusOK, responseData)
}

func GetDetailUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var user model.User
	if err := connection.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, "User not found")
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	responseData := make(map[string]interface{})
	data := make(map[string]interface{})
	data["id"] = user.Id
	data["role"] = user.Role
	data["name"] = user.Name
	data["email"] = user.Email

	responseData["data"] = data
	responseData["message"] = "Success Get Detail Users"

	helper.ResponseJson(w, http.StatusOK, responseData)
}
