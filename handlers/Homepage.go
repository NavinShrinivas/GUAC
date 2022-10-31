package handlers

import (
	"log"
	"net/http"
)

type HomepageResponse struct{
  Status bool `json:"status"`
  Message string `json:"message"`
}

func HomePage(w http.ResponseWriter, r *http.Request){
  log.Println("[ENDPOINT] Homepage")

  msg := HomePageMessage()
 	SimpleSuccessStatus(msg, w) 
  return
}

func HomePageMessage() string{
  res_msg := "You have reached the Home Page of your GUAC instance!"

  //To disaply some basic short metrics

  return res_msg
}
