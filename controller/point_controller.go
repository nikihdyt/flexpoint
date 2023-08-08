package controller

import (
	"flexpoint/connection"
	"flexpoint/helper"
	"flexpoint/model"
	"math"
	"net/http"
	"strconv"
)

func GetListPoint(w http.ResponseWriter, r *http.Request) {
	var points []model.Point

	query := r.URL.Query()
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

	var totalCount int64
	if err := db.Model(&model.Point{}).Count(&totalCount).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var totalPages int64
	if totalCount > 0 {
		totalPages = int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	} else {
		totalPages = 0
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&points).Error; err != nil {
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
	for _, point := range points {
		item := struct {
			ID    int    `json:"id"`
			UserId  int `json:"user_id"`
			Points  int `json:"points"`
		}{
			ID:    point.Id,
			UserId:  point.Id,
			Points:  point.Amount,
		}
		data = append(data, item)
	}

	responseData["data"] = data
	responseData["meta"] = meta
	responseData["message"] = "Success Get All User's Points"

	helper.ResponseJson(w, http.StatusOK, responseData)
}

