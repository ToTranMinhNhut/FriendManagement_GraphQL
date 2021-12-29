package controllers

import (
	"encoding/json"
	"net/http"
)

func MsgOK() map[string]interface{} {
	return map[string]interface{}{"success": true}
}

func MsgError(err error) map[string]interface{} {
	return map[string]interface{}{"message": err.Error(), "success": false}
}

func Message(status bool, msg string) map[string]interface{} {
	return map[string]interface{}{"message": msg, "success": status}
}

func MsgGetFriendsOk(friends []string, count int) interface{} {
	return map[string]interface{}{"count": count, "friends": friends, "success": true}
}

func MsgGetEmailReceiversOk(emails []string) interface{} {
	return map[string]interface{}{"recipients": emails, "success": true}
}

func MsgGetAllUsersOk(users []string, count int) interface{} {
	return map[string]interface{}{"count": count, "users": users, "success": true}
}

func Respond(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(statusCode)
	w.Write(response)
}
