package handlers

import (
	"GUAC/globals"
	"encoding/json"
	"log"
	"net/http"
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
