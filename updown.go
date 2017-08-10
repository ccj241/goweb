package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Izadata struct {
	Id            int
	CharacterName string
	CreateDate    string
	LastLogin     string
	Days          int
	Session       string
}

const Mysqlconf = "username:password@tcp(127.0.0.1:3306)/test"

var Iza Izadata

func UploadDatas(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("UploadDatas.gtpl")
		t.Execute(w, nil)
	} else {
		Iza.CharacterName = r.FormValue("CharacterName")
		Iza.Session = r.FormValue("Session")
		Iza.CreateDate = time.Now().Format("20060102")
		Iza.LastLogin = time.Now().AddDate(0, 0, -1).Format("20060102")
		fmt.Printf("Character Name is %v \n", Iza.CharacterName)
		fmt.Printf("Session is %v \n", Iza.Session)
		fmt.Printf("Creat Date is %v \n", Iza.CreateDate)
		fmt.Printf("LastLogin Date is %v \n", Iza.LastLogin)
		db, err := sql.Open("mysql", Mysqlconf)
		defer db.Close()
		if err != nil {
			fmt.Println("error1", err)
		}
		stmt, _ := db.Prepare("INSERT CharacterInfo SET CharacterName=?,CreatDate=?,LastLogin=?,Days=?,Session=?")
		res, error := stmt.Exec(Iza.CharacterName, Iza.CreateDate, Iza.LastLogin, 1, Iza.Session)
		if error != nil {
			fmt.Println("error3", error)
		}
		fmt.Println(res.LastInsertId())

	}
}

func GetDatas(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db, err := sql.Open("mysql", Mysqlconf)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		err := db.QueryRow("SELECT Id,Session FROM CharacterInfo WHERE LastLogin < ? limit 1", time.Now().AddDate(0, 0, -1).Format("20060102")).Scan(&Iza.Id, &Iza.Session)
		if err != nil {
			fmt.Println(err)
		}
		stmt, _ := db.Prepare("UPDATE CharacterInfo set LastLogin=? WHERE Id = ?")
		res, err := stmt.Exec(time.Now().Format("20060102"), Iza.Id)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := res.RowsAffected()
		if id == 1 {
			fmt.Println("session ", Iza.Session)
		}

	}
}

func main() {
	http.HandleFunc("/UploadDatas", UploadDatas)
	http.HandleFunc("/GetDatas", GetDatas)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
