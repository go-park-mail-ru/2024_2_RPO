package users

import (
	"RPO_back/database"
	"fmt"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hellosdfdsf world")
	database.HelloDatabase()
}
