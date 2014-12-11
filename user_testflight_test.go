package main

import (
    "github.com/zenazn/goji"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    "github.com/drewolson/testflight"
    "github.com/stretchr/testify/assert"
    "gopkg.in/yaml.v2"

    "models"
    "testing"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

func TestUserIndex(t *testing.T) {
    testflight.WithServer(rooter(goji.DefaultMux), func(r *testflight.Requester) {
        req, _ := http.NewRequest("GET", "/user/index", nil)
        req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
        response := r.Do(req)
        assert.Equal(t, 200, response.StatusCode)
    })
}

func TestUserGet(t *testing.T) {
    testflight.WithServer(rooter(goji.DefaultMux), func(r *testflight.Requester) {
        req, _ := http.NewRequest("GET", "/user/edit/22", nil)
        req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
        response := r.Do(req)
        assert.Equal(t, 200, response.StatusCode)
    })
}

func TestUserCreate(t *testing.T) {
    testflight.WithServer(rooter(goji.DefaultMux), func(r *testflight.Requester) {
        count_before := 0
        count_after := 0
        db.Table("users").Count(&count_before)

        values := url.Values{}
        values.Add("Name","testman")

        req, _ := http.NewRequest("POST", "/user/new", strings.NewReader(values.Encode()))
        req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        response := r.Do(req)

        db.Table("users").Count(&count_after)
        assert.Equal(t, 301, response.StatusCode)
        assert.Equal(t, count_before + 1, count_after)
    })
}

func TestUserCreateError(t *testing.T) {
    testflight.WithServer(rooter(goji.DefaultMux), func(r *testflight.Requester) {
        count_before := 0
        count_after := 0
        db.Table("users").Count(&count_before)

        values := url.Values{}
        values.Add("Name","エラー")

        req, _ := http.NewRequest("POST", "/user/new", strings.NewReader(values.Encode()))
        req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        response := r.Do(req)


        db.Table("users").Count(&count_after)
        assert.Equal(t, 200, response.StatusCode)
        assert.Equal(t, count_before, count_after)
    })
}

func TestUserDelete(t *testing.T) {
    testflight.WithServer(rooter(goji.DefaultMux), func(r *testflight.Requester) {
        Users := [] models.User{}
        count_before := 0
        count_after := 0
        db.Find(&Users).Count(&count_before)

        req, _ := http.NewRequest("GET", "/user/delete/1", nil)
        req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
        r.Do(req)

        db.Find(&Users).Count(&count_after)
        assert.Equal(t, count_before - 1, count_after)
    })
}

func init(){
    yml, _ := ioutil.ReadFile("conf/db.yml")
    t := make(map[interface{}]interface{})

    _ = yaml.Unmarshal([]byte(yml), &t)

    conn := t["test"].(map[interface {}]interface {})
    db, _ = gorm.Open("mysql", conn["user"].(string)+conn["password"].(string)+"@/"+conn["db"].(string)+"?charset=utf8&parseTime=True")
    db.DropTable(&models.User{})
    db.CreateTable(&models.User{})
    User := models.User{Name: "deluser"}
    db.Save(&User)    
}
