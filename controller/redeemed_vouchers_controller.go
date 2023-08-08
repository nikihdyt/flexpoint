package controller

import (
	"encoding/json"
	"flexpoint/connection"
	"flexpoint/helper"
	"flexpoint/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RedeemVoucher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var voucher map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&voucher); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	pointUsed := int(voucher["point_used"].(float64))

	userCurrentPoints, err := GetUserCurrentPoints(int(userID))
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userCurrentPoints < pointUsed {
		ResponseError(w, http.StatusBadRequest, "Point belum cukup. Kumpulkan lebih banyak point.")
		return
	}

	updatedPoints := userCurrentPoints - pointUsed
	log.Println("updatedPoints: ", updatedPoints)
	if err := connection.DB.Model(&model.Point{}).Where("user_id = ?", userID).Update("amount", updatedPoints).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	query := `INSERT INTO redeemed_vouchers (user_id, voucher_code, redeemed_date, point_used)
			  VALUES (?, ?, ?, ?)`
	if err := connection.DB.Exec(query, userID, voucher["voucher_code"], voucher["redeemed_date"], voucher["point_used"]).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := make(map[string]interface{})
	responseData["message"] = "Voucher redeemed successfully. Your current points: " + strconv.Itoa(updatedPoints)

	helper.ResponseJson(w, http.StatusOK, responseData)

}

func GetUserCurrentPoints(userID int) (int, error) {
	var userPoint model.Point
	if err := connection.DB.Where("user_id = ?", userID).First(&userPoint).Error; err != nil {
		if err != nil {
			return 0, err
		}
	}
	return userPoint.Amount, nil
}
