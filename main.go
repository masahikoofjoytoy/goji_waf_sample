package main

import (
    "github.com/zenazn/goji"
    "github.com/jinzhu/gorm"
    "github.com/zenazn/goji/web"
    _ "github.com/go-sql-driver/mysql"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
    "net/http"
)

var db  gorm.DB

func main() {
    connect()
    rooter(goji.DefaultMux)
    goji.Serve()
}

func rooter(m *web.Mux) http.Handler {
    m.Use(SuperSecure)
    m.Get("/index", UserRoot)
    m.Get("/user/index", UserIndex)
    m.Get("/user/new", UserNew)
    m.Post("/user/new", UserCreate)
    m.Get("/user/edit/:id", UserEdit)
    m.Post("/user/update/:id", UserUpdate)
    m.Get("/user/delete/:id", UserDelete)
    
    return m
}

func connect(){
    yml, err := ioutil.ReadFile("conf/db.yml")
    if err != nil {
        panic(err)
    }

    t := make(map[interface{}]interface{})

    _ = yaml.Unmarshal([]byte(yml), &t)

    conn := t[os.Getenv("GOJIENV")].(map[interface {}]interface {})
    db, err = gorm.Open("mysql", conn["user"].(string)+conn["password"].(string)+"@/"+conn["db"].(string)+"?charset=utf8&parseTime=True")
    if err != nil {
        panic(err)
    }
}
