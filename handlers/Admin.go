package handlers

import (
	// "log"
	"GUAC/globals"
)

func CheckAdminExists(admin_id string) bool {
	constraint := globals.Admins{
		Adm_id: admin_id,
	}
	var results []globals.Admins
	globals.DbConn.Where(&constraint).Find(&results)

	for _, v := range results {
		if v.Adm_id == admin_id {
			return true
		}
	}
	return false

}
