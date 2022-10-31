package main

import(
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "log"
  "net/http"
  "GUAC/handlers"
  "GUAC/globals"
  "sync"
)




func RoutesFunction(){
  http.HandleFunc("/", handlers.HomePage)
  http.HandleFunc("/test", handlers.BasicTest)
  log.Fatal(http.ListenAndServe("0.0.0.0:3030",nil))
}

func DatabaseEnv() *gorm.DB{
  log.Println("[DATABASE] Intiating connection to MYSQL instance.")
  db_conf := mysql.Config{
    //[TODO]to convert connection string to env file.
    DSN: "common_user:common_pass@tcp(127.0.0.1:3306)/testdbrun?charset=utf8mb4&parseTime=True&loc=Local", // data source name
    DefaultStringSize: 256, // default size for string fields
    DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
    DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
    DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
    SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
  }
  db_con,err := gorm.Open(mysql.New(db_conf))
  if err!=nil{
    //GORM prints the error
    log.Fatal("[DATABASE] Error accessing databse, please ensure you have run our setup files right. Also do make sure MYSQL service is running.")
  }
  log.Println("[DATABASE] Connection successfull. Migrating/Updating schemas. Lossless update.")
  db_con.AutoMigrate(
    &globals.Admins{},
    &globals.Doc{},
    &globals.User_perms{},
    &globals.Auth_code{},
  )
  log.Println("[DATABASE] Migrations completete.")
  return db_con
}

func main(){
  wg := new(sync.WaitGroup)
  log.Println("[DEBUG] Starting server on port 3030")
  globals.DbConn = DatabaseEnv()
  defer log.Println("[DATABASE] Intial database environment complete!")
  wg.Add(1)
  go RoutesFunction()
  log.Println("[DEBUG] Server listening on port 3030")
  wg.Wait()
}
