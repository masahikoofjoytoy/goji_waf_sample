package main

import (
    "github.com/zenazn/goji/web"
    "github.com/wcl48/valval"

    "net/http"
    "encoding/base64"
    "strings"
    "html/template"
    "models"
    "fmt"
    "strconv"
)

var tpl *template.Template
const Password = "user:user"

type FormData struct{
    User models.User
    Mess string
}

func UserRoot(c web.C, w http.ResponseWriter, r *http.Request) {
    fmt.Println("UserRoot")
    http.Redirect(w, r, "/user/index", 301)
 }

func UserIndex(c web.C, w http.ResponseWriter, r *http.Request) {
    Users := [] models.User{}
    db.Find(&Users)
    tpl = template.Must(template.ParseFiles("view/user/index.html"))
    tpl.Execute(w,Users)
}

func UserNew(c web.C, w http.ResponseWriter, r *http.Request){
    tpl = template.Must(template.ParseFiles("view/user/new.html"))
    tpl.Execute(w,FormData{models.User{}, ""})    
}

func UserCreate(c web.C, w http.ResponseWriter, r *http.Request){
    User := models.User{Name: r.FormValue("Name")}
    if err := models.UserValidate(User); err != nil {
        var Mess string
        errs := valval.Errors(err)
        for _, errInfo := range errs {
            Mess += fmt.Sprint(errInfo.Error)
        }
        tpl = template.Must(template.ParseFiles("view/user/new.html"))
        tpl.Execute(w,FormData{User, Mess})    
    } else {
        db.Create(&User)    
        http.Redirect(w, r, "/user/index", 301)
    }
}

func UserEdit(c web.C, w http.ResponseWriter, r *http.Request){
    User := models.User{}
    User.Id, _ = strconv.ParseInt(c.URLParams["id"], 10, 64)
    db.Find(&User)
    tpl = template.Must(template.ParseFiles("view/user/edit.html"))
    tpl.Execute(w,FormData{User, ""})    
}

func UserUpdate(c web.C, w http.ResponseWriter, r *http.Request){
    User := models.User{}
    User.Id, _ = strconv.ParseInt(c.URLParams["id"], 10, 64)
    db.Find(&User)
    User.Name = r.FormValue("Name")
    if err := models.UserValidate(User); err != nil {
        var Mess string
        errs := valval.Errors(err)
        for _, errInfo := range errs {
            Mess += fmt.Sprint(errInfo.Error)
        }
        tpl = template.Must(template.ParseFiles("view/user/edit.html"))
        tpl.Execute(w,FormData{User, Mess})
    } else {
        db.Save(&User)    
        http.Redirect(w, r, "/user/index", 301)
    }
}

func UserDelete(c web.C, w http.ResponseWriter, r *http.Request){
    User := models.User{}
    User.Id, _ = strconv.ParseInt(c.URLParams["id"], 10, 64)
    db.Delete(&User)
    http.Redirect(w, r, "/user/index", 301)        
}

func SuperSecure(c *web.C, h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Basic ") {
            pleaseAuth(w)
            return
        }

        password, err := base64.StdEncoding.DecodeString(auth[6:])
        if err != nil || string(password) != Password {
            pleaseAuth(w)
            return
        }

        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

func pleaseAuth(w http.ResponseWriter) {
    w.Header().Set("WWW-Authenticate", `Basic realm="Gritter"`)
    w.WriteHeader(http.StatusUnauthorized)
    w.Write([]byte("Go away!\n"))
}
