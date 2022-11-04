package handlers

import (
	// "log"
	"GUAC/globals"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

type AdminRequest struct {
	Adm_id             string `json:"adm_id"`
	Adm_pass_plaintext string `json:"adm_pass_plaintext"`
}

type RKGSReqRes struct {
	Work  string `json:"Work"`
	Pool  string `json:"Pool"`
	Len   string `json:"Len"`
	Url   string `json:"Url"`
	Error string `json:"Error"`
}

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

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("[ENDPOINT] Creat Admin endpoint accessed.")
	if r.TLS == nil {
		//TLS attribute in request is filled server side, hence this comparision in secure and cannot be manipulated!
		//TLS was not used for this request
		log.Println("[DEBUG] Unauthorised access through non TLS port.")
		SimpleFailStatus("Please use secure Socket port 3031 for admins and auth.", w)
		return
	} else {
		//[TODO] Salting, peperring and Hash needs to be done.
		//Until then passwords are stored in plain text

		request_body := AdminRequest{}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&request_body)

		if err != nil {
			log.Println("[ENDPOINT ERROR] Error reading body, possibly malformed inputs or wrong fields.")
			SimpleFailStatus("Please give proper inputs!", w)
			return
		}
		if CheckAdminExists(request_body.Adm_id) {
			SimpleFailStatus("Admin with that user name already exists :(", w)
			return
		} else {
			//Create admin
			insert_admin_record := globals.Admins{
				Adm_id:        request_body.Adm_id,             //Need to sanitize user names
				Adm_hash_pass: request_body.Adm_pass_plaintext, //[TODO] Salting Hashing and peperring
			}
			log.Println("[DATABASE] Creating admin record in database")
			result := globals.DbConn.Create(&insert_admin_record)
			if result.Error != nil {
				log.Println("[DATABASE] Something went wrong, Admin.go", err)
				SimpleFailStatus("We are facing error on our side, please try again later.", w)
				return
			}
			SimpleSuccessStatus("Succesfully created Admin account!", w)
		}
	}
}

func GetAuthCode(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		//TLS attribute in request is filled server side, hence this comparision in secure and cannot be manipulated!
		//TLS was not used for this request
		log.Println("[DEBUG] Unauthorised access through non TLS port.")
		SimpleFailStatus("Please use secure Socket port 3031 for admins and auth.", w)
		return
	}

	//Need to request auth code from RKGS
	RKGSres := RKGSReqRes{
		Work:  "generate",
		Pool:  "GUAC_pool",
		Url:   "",
		Error: "",
	}
	RKGSserverAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:6379")
	if err != nil {
		log.Println("[ERROR] Make sure RKGS server is listening on port 6379")
	}
	RKGSclientAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	request_json, err := json.Marshal(RKGSres)

	if err != nil {
		log.Println("[ERROR] Error creating json payload for RKGS")
	}

}
