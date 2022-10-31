package handlers

import (
	"GUAC/globals"
	"log"
	"net/http"
)

type ApitestResponse struct{
  Status bool `json:"status"`
  Test_record string `json:"test_record"`
}

func BasicTest(w http.ResponseWriter, r *http.Request){
  log.Println("[ENDPOINT] Api test...Running basic test on DB conn.")
  test_record := globals.Admins{
    Adm_id : "test12345",
    Adm_hash_pass: "hash_pass_test!@##%@",
  }
  test_record2 := globals.Admins{
    Adm_id : "test12345",
    Adm_hash_pass: "hash_pass_test!@##%@",
  }
  var result []globals.Admins

  log.Println("[DATABASE] Inserting test records.")
  globals.DbConn.Create(&test_record)
  globals.DbConn.Create(&test_record2)

	println()
  log.Println("[DATABASE] Retrieving test records and returning to API.")
  conditions := globals.Admins{
    Adm_id: "test12345",
  }
	msg := "" 
  globals.DbConn.Where(&conditions).Find(&result)
  for _,v := range result {
    msg += v.Adm_id+"   "+v.Adm_hash_pass+"   "
  }
  SimpleSuccessStatus(msg, w)

  log.Println("[DATABASE] Deleting test record from database")
  globals.DbConn.Delete(&conditions)
  return
}

