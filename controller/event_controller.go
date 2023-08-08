package controller

import (
	"encoding/json"
	"flexpoint/connection"
	"flexpoint/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {

	var event model.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	if err := connection.DB.Create(&event).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := make(map[string]interface{})
	responseData["message"] = "Success Create Event"

	ResponseJson(w, http.StatusCreated, responseData)

}

func VerifyEventAndAddPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var points map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&points); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	var event model.Event
	if err := connection.DB.First(&event, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, "Event not found")
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if event.IsVerified {
		ResponseError(w, http.StatusBadRequest, "Event already verified")
		return
	}

	if err := connection.DB.Model(&event).Where("id = ?", id).Update("is_verified", true).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: add points to user
	var point model.Point
	point.UserId = event.UserId
	point.Amount = int(points["points"].(float64))
	// find row with user_id = event.UserId
	var pointExist model.Point
	if err := connection.DB.Where("user_id = ?", event.UserId).First(&pointExist).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			if err := connection.DB.Create(&point).Error; err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
		default:
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// update row with user_id = event.UserId
		if err := connection.DB.Model(&pointExist).Where("user_id = ?", event.UserId).Update("amount", pointExist.Amount+point.Amount).Error; err != nil {
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	responseData := make(map[string]interface{})
	responseData["message"] = "Success Verify Event"

	ResponseJson(w, http.StatusOK, responseData)

}

func GetEventDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var event model.Event
	if err := connection.DB.First(&event, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, "Event not found")
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	responseData := make(map[string]interface{})
	data := make(map[string]interface{})
	data["id"] = event.Id
	data["user_id"] = event.UserId
	data["event_name"] = event.EventName
	data["event_date"] = event.EventDate
	data["url"] = event.URL
	data["is_verified"] = event.IsVerified
	data["points"] = event.Points

	responseData["data"] = data
	responseData["message"] = "Success Get Detail Event"

	ResponseJson(w, http.StatusOK, responseData)
}
