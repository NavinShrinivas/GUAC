package handlers

import (
	"GUAC/globals"
	"encoding/json"
	"log"
	"net/http"
	"errors"
)

type DocInfoRequest struct {
	Doc_id string `json:"doc_id"`
}

type DocInfoCustomReponse struct {
	Doc_id      string `json:"doc_id"`
	Def_permbit int8   `json:"def_permbit"`
	Adm_id      string `json:"adm_id"`
}

func ReturnDocInfo(w http.ResponseWriter, v globals.Doc) {
	DocInfoCustomReponse := DocInfoCustomReponse{
		Doc_id:      v.Doc_id,
		Def_permbit: v.Def_permbit,
		Adm_id:      v.Adm_id,
	}
	response_obj, err := json.Marshal(DocInfoCustomReponse)
	if err != nil {
		log.Println("[ERROR] Error Parsing struct in Document end point.")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response_obj)
	return
}

func DocInfo(w http.ResponseWriter, r *http.Request) {

	log.Println("[ENDPOINT] Document GET")
	request_body := DocInfoRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request_body)

	if err != nil {
		SimpleFailStatus("Something went wrong on our side, plase try again later!", w)
		return
	}
	//else :
	var db_res []globals.Doc
	constraints := globals.Doc{
		Doc_id: request_body.Doc_id,
	}
	globals.DbConn.Where(&constraints).Find(&db_res)
	if len(db_res) == 0 {
		SimpleFailStatus("No document with that ID was found!", w)
	}
	for _, v := range db_res {
		if v.Doc_id == request_body.Doc_id {
			ReturnDocInfo(w, v)
			return
		}
	}
}


type InsertdocRequest struct {
	Doc_id      string `json:"doc_id"`
	Adm_id      string `json:"adm_id"`
	Def_permbit string `json:"def_permbit"`
}

func bitstringtoint8(bitstring string) (int8, error) {
	if len(bitstring) > 8 {
		return int8(0), errors.New("Bit string too long")
	}
	ans := int8(0)
	for _, v := range bitstring {
		if v == '1' {
			ans = ans<<1 | 1
		} else if v == '0' {
			ans = ans << 1
		} else {
			return int8(0), errors.New("Bit string malformed!")
		}
	}
	return ans, nil
}

func Insertdoc(w http.ResponseWriter, r *http.Request) {
	log.Println("[ENDPOINT] Inserting new document record.")

	request_body := InsertdocRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&request_body)

	if err != nil {
		log.Println("[ENDPOINT ERROR] Error reading body, possibly malformed inputs or wrong fields.")
		SimpleFailStatus("Please give proper inputs!", w)
		return
	}

	//Need to check if admin account exists, else request to create admin.
	if CheckAdminExists(request_body.Adm_id) {
		inserting_permbit, err := bitstringtoint8(request_body.Def_permbit)
		if err != nil {
			log.Println("[ENDPOINT ERROR] Recieved invalid perm bit")
			SimpleFailStatus("Invalid perm bit string!", w)
			return
		}
		new_record := globals.Doc{
			Doc_id:      request_body.Doc_id,
			Adm_id:      request_body.Adm_id,
			Def_permbit: inserting_permbit,
		}
		insert_err := globals.DbConn.Create(&new_record)
		if insert_err.Error != nil {
			SimpleFailStatus(insert_err.Error.Error(), w)
			return
		}
		SimpleSuccessStatus("Inserted document, document now ready to track permissions!", w)
		return
	} else {
		log.Println("[ENDPOINT ERROR] Document trying to be created by unavailable admin!")
		SimpleFailStatus("Admin doesn't exist.", w)
	}
}
