package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type HomepageResponse struct{
  Status bool `json:"status"`
  Message string `json:"message"`
}

func HomePage(w http.ResponseWriter, r *http.Request){
  log.Println("[ENDPOINT] Homepage")

  response := HomepageResponse{
    Status : true,
    Message: HomePageMessage(),
  }

  w.WriteHeader(http.StatusOK)
  response_json,err := json.Marshal(response)

  if err!=nil{
    log.Println("[error] error marshallling structure to json. err id : 1")
  }

  w.Write(response_json)
  return
}

func HomePageMessage() string{
  res_msg := "You have reached the Home Page of your GUAC instance!"

  //To disaply some basic short metrics

  return res_msg
}
