package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

func PingDB(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDb, err := db.DB()
		if err != nil || sqlDb.Ping() != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
