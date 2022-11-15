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

type AdminAuthCodeResponse struct {
	Adm_id    string `json:"adm_id"`
	Auth_Code string `json:"auth_code"`
}

type RKGSReqRes struct {
	Work  string `json:"Work"`
	Pool  string `json:"Pool"`
	Len   string `json:"Len"`
	Url   string `json:"Url"`
	Error string `json:"Error"`
}

func AuthCodeValid(adm_id string, auth_code string) (bool, error) {

	constraint := globals.Auth_code{
		Adm_id: adm_id,
	}
	var db_res []globals.Auth_code
	globals.DbConn.Where(&constraint).Find(&db_res)

	for _, v := range db_res {
		if v.Adm_id == adm_id && v.Auth_code == auth_code {
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
	log.Println("[ENDPOINT] Create Admin endpoint accessed.")
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

	//First check if admin exist, if so verify password
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
		//Verify password [TODO] : Need this part of code to be extremely secure.
		if VerifyAdminPassword(request_body) {
			auth_code, err := get_auth_code_from_rkgs(w)
			if err != nil {
				return
			}
			//Else, we have to insert auth code to database and return it
			db_new_record := globals.Auth_code{
				Adm_id:    request_body.Adm_id,
				Auth_code: auth_code,
			}
			result := globals.DbConn.Create(db_new_record)
			if result.Error != nil {
				//Automatically log prints
				SimpleFailStatus("Could not generate auth code, please try again later.", w)
				return
			}
			res := AdminAuthCodeResponse{
				Adm_id:    db_new_record.Adm_id,
				Auth_Code: db_new_record.Auth_code,
			}
			res_json, err := json.Marshal(res)
			if err != nil {
				log.Println("[ENDPOINT ERROR] Error send back result for auth code")
				SimpleFailStatus("Please try again, we ran into a small error on our side :(", w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(res_json)
			return
		} else {
			SimpleFailStatus("Wrong credentials!!", w)
			return
		}
	} else {
		SimpleFailStatus("Admin with that user name does not exist, please create admin account before trying actions.", w)
		return
	}

}

func VerifyAdminPassword(request_body AdminRequest) bool {
	db_record_constraint := globals.Admins{
		Adm_id: request_body.Adm_id,
	}
	db_record_response := globals.Admins{}
	globals.DbConn.First(&db_record_response).Where(&db_record_constraint)
	//[TODO] Need to calculate password hashes here once we start storing hashed pass in database
	if request_body.Adm_pass_plaintext == db_record_response.Adm_hash_pass {
		return true
	}
	return false
}

func get_auth_code_from_rkgs(w http.ResponseWriter) (string, error) {
	//Need to request auth code from RKGS
	RKGSreq := RKGSReqRes{
		Work:  "generate",
		Pool:  "GUAC_pool",
		Len:   "64",
		Url:   "",
		Error: "",
	}
	RKGSserverAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:5001")
	if err != nil {
		log.Println("[ERROR] Make sure RKGS server is listening on port 5001")
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		return "", err
	}
	RKGSclientAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	request_json, err := json.Marshal(RKGSreq)
	if err != nil {
		log.Println("[ERROR] Error creating json payload for RKGS")
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		return "", err
	}
	RKGSconn, err := net.DialUDP("udp", RKGSclientAddr, RKGSserverAddr)
	if err != nil {
		log.Println("[ERROR] Make sure RKGS server is running on port 5001")
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		return "", err
	}
	//First try sending payload to server :
	_, write_err := RKGSconn.Write(request_json)
	if write_err != nil {
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		log.Fatal("[ERROR] Something is goig quite wrong, error writing to RKGS server.")
	}
	received := make([]byte, 2048)
	bits, err := RKGSconn.Read(received)
	if err != nil {
		log.Println("[ERROR] Make sure RKGS server is running on port 5001")
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		return "", err
	}
	RKGSres := RKGSReqRes{}
	err = json.Unmarshal(received[:bits], &RKGSres)
	if err != nil {
		log.Println("[ERROR] Error Unmarshsalling reponse from RKGS")
		SimpleFailStatus("We are facing some errors on our side, please try again later", w)
		return "", err
	}
	return RKGSres.Url, nil
}
