package handlers

import (
	"GUAC/globals"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

type UserCreationRequest struct {
	User_id    string `json:"user_id"`
	Doc_id     string `json:"doc_id"`
	Adm_id     string `json:"adm_id"`
	Auth_code  string `json:"auth_code"`
	Nd_permbit string `json:"nd_permbit"`
}


func Users(w http.ResponseWriter, r *http.Request) {
	log.Println("[ENDPOINT] Access on Users endpoint")

	if r.Method == "POST" {
		//Validate Auth code
		request_body := UserCreationRequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&request_body)
		valid, err := AuthCodeValid(request_body.Adm_id, request_body.Auth_code)

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
	if r.Method == "PUT"{
		UserModify(w,r)
		return
	}
	if r.Method == "DELETE"{
		// New funtion to allow deletion after veriyfinh auth code
		return
	}
}

type UserModifyRequest struct{
	User_id    string `json:"user_id"`
	Doc_id     string `json:"doc_id"`
	Adm_id     string `json:"adm_id"`
	Auth_code  string `json:"auth_code"`
	New_Nd_permbit string `json:"new_nd_permbit"`
}

func UserModify(w http.ResponseWriter, r *http.Request){
		request_body := UserModifyRequest{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&request_body)
		valid, err := AuthCodeValid(request_body.Adm_id, request_body.Auth_code)

		if err != nil {
			SimpleFailStatus("Something is going very wrong, Please try again later.", w)
			return
		}
		if valid {
			//Time to modify user
			int8_bit_string, err := bitstringtoint8(request_body.New_Nd_permbit)
			if err != nil {
				SimpleFailStatus("Invalid perm bit", w)
				return
			}
			db_old_record := globals.User_perms{
				Doc_id:     request_body.Doc_id,
			}
			fmt.Println(int8_bit_string, request_body.New_Nd_permbit)
			result := globals.DbConn.Model(&db_old_record).Where("user_id",request_body.User_id).Where("doc_id", request_body.Doc_id).Update("nd_permbit",int8_bit_string)
			if result.RowsAffected == 0{
				log.Println("[ENDPOINT] No record with that match found!")
				SimpleFailStatus("Nice Try, that user does not exists lol!",w)
				return
			}
			if result.Error != nil {
				log.Println("[ENDPOINT ERROR] Error Updating user record!!!!")
				SimpleFailStatus("We ran into an error updating user records, please try again later.", w)
				return
			}
			SimpleSuccessStatus("User record for asked document updated successfully!", w)
		} else {
			SimpleFailStatus("Invalid Auth code or Admin. Please verify your json inputs again.", w)
			return
		}
}
