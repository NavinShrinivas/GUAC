package main

import (
	"GUAC/globals"
	"GUAC/handlers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

// ----------Main network functions----------
// Creating method splitter, Not including in handlerFuncs
func documentRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//Creating new document entry flow
		handlers.Insertdoc(w, r)
	} else if r.Method == "GET" {
		//Get info about an existing document entry
		handlers.DocInfo(w, r)
	}
}

func adminRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//Flow to create new admin
		handlers.CreateAdmin(w, r)
	} else if r.Method == "GET" {
		//Get info about admins existance
	}
}

func RoutesFunction() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/test", handlers.BasicTest)
	http.HandleFunc("/document", documentRouter)
	http.HandleFunc("/admin", adminRouter)
	http.HandleFunc("/authcode", handlers.GetAuthCode)
	http.HandleFunc("/users", handlers.Users)
	http.HandleFunc("/access", handlers.Access)
	log.Fatal(http.ListenAndServe("0.0.0.0:3030", nil))
}

var certFile = "./localhost.crt" //Need to changed if we enter product build
var keyFile = "./localhost.key"  //Need to changed if we enter product build
func SecureServerInit() {
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:3031", certFile, keyFile, nil))
	//This server also serves same path as unsecure port, just complies to TLS
}

//End of main network funcs [Handlers, Method diff, Server]

// ----------main runnner---------
func main() {
	wg := new(sync.WaitGroup)
	log.Println("[DEBUG] Starting server on port 3030")
	globals.DbConn = DatabaseEnv()
	defer log.Println("[DATABASE] Intial database environment complete!")
	wg.Add(1)
	go RoutesFunction()
	go SecureServerInit()
	log.Println("[DEBUG] Server listening on port 3030")
	log.Println("[DEBUG] Server listening on port 3031 For secure admin transactions")
	wg.Wait()
}

//main function ends

//----------Databse env setup----------

func DatabaseEnv() *gorm.DB {
	log.Println("[DATABASE] Intiating connection to MYSQL instance.")
	db_conf := mysql.Config{
		//[TODO]to convert connection string to env file.
		DSN:                       "common_user:common_pass@tcp(127.0.0.1:3306)/testdbrun?charset=utf8mb4&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                              // default size for string fields
		DisableDatetimePrecision:  true,                                                                                             // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                             // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                             // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                            // auto configure based on currently MySQL version
	}
	db_con, err := gorm.Open(mysql.New(db_conf))
	if err != nil {
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

//End of Main database operations
