package handlers

import (
	// "encoding/json"
	"GUAC/globals"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type InsertdocRequest struct{
  Doc_id string  `json:"doc_id"`
  Adm_id string  `json:"adm_id"`
  Def_permbit string `json:"def_permbit"`
}


func bitstringtoint8(bitstring string) (int8, error){
  if len(bitstring) > 8{
    return int8(0),errors.New("Bit string too long")
  }
  ans := int8(0)
  for _,v := range bitstring{
    if v == '1'{
      ans = ans<< 1 | 1
    }else if v == '0'{
      ans = ans << 1
    }else{
      return int8(0),errors.New("Bit string malformed!")
    }
  } 
  return ans,nil
}

func Insertdoc(w http.ResponseWriter, r *http.Request){
  log.Println("[ENDPOINT] Inserting new document record.")


  request_body := InsertdocRequest{}
  dec := json.NewDecoder(r.Body)
  dec.DisallowUnknownFields()
  err:= dec.Decode(&request_body)

  if err!=nil{
    log.Println("[ENDPOINT ERROR] Error reading body, possibly malformed inputs or wrong fields.")
    SimpleFailStatus("Please give proper inputs!", w)
    return
  }

  //Need to check if admin account exists, else request to create admin.
  if CheckAdminExists(request_body.Adm_id){
    inserting_permbit, err := bitstringtoint8(request_body.Def_permbit)
    if err!=nil {
      log.Println("[ENDPOINT ERROR] Recieved invalid perm bit")
      SimpleFailStatus("Invalid perm bit string!", w)
      return
    }
    new_record := globals.Doc{
      Doc_id: request_body.Doc_id,
      Adm_id: request_body.Adm_id,
      Def_permbit : inserting_permbit,
    }
    insert_err := globals.DbConn.Create(&new_record)
		if insert_err.Error!=nil{
			SimpleFailStatus(insert_err.Error.Error(), w)
			return
		}
    SimpleSuccessStatus("Inserted document, document now ready to track permissions!", w)
    return
  }else{
    log.Println("[ENDPOINT ERROR] Document trying to be created by unavailable admin!")
    SimpleFailStatus("Admin doesn't exist.", w)
  }
}

