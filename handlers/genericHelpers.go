package handlers
import(
  "net/http"
  "encoding/json"
  "log"
)
type Response struct{
  Status bool `json:"status"`
  Message string `json:"message"`
}

func SimpleFailStatus(msg string, w http.ResponseWriter){
  failed_reponse := Response{
    Status: false,
    Message: msg,
  }
  failed_reponse_json, err := json.Marshal(failed_reponse)
  if err!=nil{
    log.Println("[ERROR] error marshallling structure to json.")
    return
  }
  w.WriteHeader(http.StatusForbidden)
  w.Write(failed_reponse_json)
}

func SimpleSuccessStatus(msg string, w http.ResponseWriter){
  failed_reponse := Response{
    Status: true,
    Message: msg,
  }
  failed_reponse_json, err := json.Marshal(failed_reponse)
  if err!=nil{
    log.Println("[ERROR] error marshallling structure to json.")
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(failed_reponse_json)
}

