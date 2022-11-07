package handlers

import (
	"GUAC/globals"
	"encoding/json"
	"log"
	"net/http"
)

type AccessDocRequest struct {
	User_id string `json:"user_id"`
	Doc_id  string `json:"doc_id"`
}

type AccessDocResponse struct {
	User_id    string `json:"user_id"`
	Doc_id     string `json:"doc_id"`
	Nd_permbit string `json:"nd_permbit"`
}

func int8tobitstring(num int8) string {
	ans := ""
	var k int8
	for c := 8; c >= 0; c-- {
		k = num >> c
		if k&1 == 1 {
			ans += "1"
		} else {
			ans += "0"
		}
	}
	return ans
}

func Access(w http.ResponseWriter, r *http.Request) {
	log.Println("[ENDPOINT] Access on Access endpoint")
	//We do not need Auth code or TLS for this endpoint

	request_body := AccessDocRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request_body)
	if err != nil {
		log.Println("[ERROR] Unmarshalling error")
		SimpleFailStatus("Something is going wrong in our side, please try again later", w)
		return
	}

	user_search_constraint := globals.User_perms{
		User_id: request_body.User_id,
	}

	var db_res []globals.User_perms

	globals.DbConn.Where(&user_search_constraint).Find(&db_res)

	//[TODO]If len of db_res is greater than 1, we have some junk in our systems and we need to send WARNS out

	for _, v := range db_res {
		if v.User_id == request_body.User_id && v.Doc_id == request_body.Doc_id {
			response := AccessDocResponse{
				Doc_id:     v.Doc_id,
				User_id:    v.User_id,
				Nd_permbit: int8tobitstring(v.Nd_permbit),
			}
			response_stream, err := json.Marshal(response)
			if err != nil {
				log.Println("[ERROR] Erro Marshalling reponse")
				SimpleFailStatus("Something went wrong on our end, please try again later.", w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response_stream)
			return
		}
	}
}
