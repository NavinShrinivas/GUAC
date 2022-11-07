package handlers

import (
	"GUAC/globals"
	"encoding/json"
	"log"
	"net/http"
)

type UserCreationRequest struct {
	User_id    string `json:"user_id"`
	Doc_id     string `json:"doc_id"`
	Adm_id     string `json:"adm_id"`
	Auth_code  string `json:"auth_code"`
	Nd_permbit string `json:"nd_permbit"`
}

func AuthCodeValid(request_body UserCreationRequest) (bool, error) {

	constraint := globals.Auth_code{
		Adm_id: request_body.Adm_id,
	}
	var db_res []globals.Auth_code
	globals.DbConn.Where(&constraint).Find(&db_res)

	for _, v := range db_res {
		if v.Adm_id == request_body.Adm_id && v.Auth_code == request_body.Auth_code {
			delete_constarint := globals.Auth_code{
				Adm_id:    v.Adm_id,
				Auth_code: v.Auth_code,
			}
			response := globals.DbConn.Delete(delete_constarint) //Auth code acts as a one time password
			if response.Error == nil {
				log.Println("[DATABASE] Removed Auth code after use!")
			}
			return true, nil
		}
	}

	return false, nil
}

func Users(w http.ResponseWriter, r *http.Request) {
	log.Println("[ENDPOINT] Access on Users endpoint")

	if r.Method == "POST" {
		//Validate Auth code
		request_body := UserCreationRequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&request_body)
		valid, err := AuthCodeValid(request_body)

		if err != nil {
			SimpleFailStatus("Something is going very wrong, Please try again later.", w)
			return
		}
		if valid {
			//Time to create user
			int8_bit_string, err := bitstringtoint8(request_body.Nd_permbit)
			if err != nil {
				SimpleFailStatus("Invalid perm bit", w)
				return
			}
			db_new_record := globals.User_perms{
				Doc_id:     request_body.Doc_id,
				User_id:    request_body.User_id,
				Nd_permbit: int8_bit_string,
			}
			result := globals.DbConn.Create(db_new_record)
			if result.Error != nil {
				log.Println("[ENDPOINT ERROR] Error inserting User to databse!!!!")
				SimpleFailStatus("We ran into an error inserting user records, please try again later.", w)
				return
			}
			SimpleSuccessStatus("User record for asked document inserted successfully!", w)
		} else {
			SimpleFailStatus("Invalid Auth code or Admin. Please verify your json inputs again.", w)
			return
		}
	}
	if r.Method == "GET" {
		SimpleFailStatus("Invalid METHOD for endpoint [user]", w)
		return
	}
}
