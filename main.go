package main

import(
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "log"
  "net/http"
  "GUAC/handlers"
  "sync"
)


func RoutesFunction(){
  http.HandleFunc("/", handlers.HomePage)
  log.Fatal(http.ListenAndServe("0.0.0.0:3030",nil))
}

func DatabaseEnv() /* *gorm.DB */{
  db_conf := mysql.Config{
    DSN: "common:common@tcp(127.0.0.1:3306)/permdb?charset=utf8mb4&parseTime=True&loc=Local", // data source name
    DefaultStringSize: 256, // default size for string fields
    DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
    DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
    DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
    SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
  }
  _,err := gorm.Open(mysql.New(db_conf))

  if err!=nil{
    //GORM prints the error
    log.Fatal("[DATABASE] Error accessing databse, please ensure you have run our setup files right.")
  }
}

func main(){
  wg := new(sync.WaitGroup)
  log.Println("[DEBUG] Starting server on port 3030")
  DatabaseEnv()
  wg.Add(1)
  go RoutesFunction()
  log.Println("[DEBUG] Server listening on port 3030")
  wg.Wait()
}
