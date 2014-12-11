package main

import (
    "github.com/jinzhu/gorm"
    "github.com/zenazn/goji/web"
    _ "github.com/go-sql-driver/mysql"
    "github.com/stretchr/testify/assert"
    "gopkg.in/yaml.v2"

    "models"
    "net/http"
    "io/ioutil"
    "net/url"
    "strings"
    "testing"
    "net/http/httptest"
)

func TestUserIndex(t *testing.T) { 
    m := web.New()
    rooter(m)
    ts := httptest.NewServer(m)
    defer ts.Close()

    req, _ := http.NewRequest("GET", ts.URL + "/user/index", nil)
    req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
    client := new(http.Client)
    response, _ := client.Do(req)

    assert.Equal(t, 200, response.StatusCode)
}

func TestUserGet(t *testing.T) { 
    m := web.New()
    rooter(m)
    ts := httptest.NewServer(m)
    defer ts.Close()

    req, _ := http.NewRequest("GET", ts.URL + "/user/edit/1", nil)
    req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
    client := new(http.Client)
    response, _ := client.Do(req)

    assert.Equal(t, 200, response.StatusCode)
}

func TestUserCreate(t *testing.T) {
    count_before := 0
    count_after := 0

    m := web.New()
    rooter(m)
    ts := httptest.NewServer(m)
    defer ts.Close()

    db.Table("users").Count(&count_before)
    values := url.Values{}
    values.Add("Name","testman")

    req, _ := http.NewRequest("POST", ts.URL + "/user/new",strings.NewReader(values.Encode()))
    req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := new(http.Client)
    response, _ := client.Do(req)

    db.Table("users").Count(&count_after)

    assert.Equal(t, 301, response.StatusCode)
    assert.Equal(t, count_before + 1, count_after)
}

func TestUserCreateError(t *testing.T) {
    count_before := 0
    count_after := 0

    m := web.New()
    rooter(m)
    ts := httptest.NewServer(m)
    defer ts.Close()

    db.Table("users").Count(&count_before)
    values := url.Values{}
    values.Add("Name","エラー")

    req, _ := http.NewRequest("POST", ts.URL + "/user/new",strings.NewReader(values.Encode()))
    req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := new(http.Client)
    response, _ := client.Do(req)

    db.Table("users").Count(&count_after)

    assert.Equal(t, 200, response.StatusCode)
    assert.Equal(t, count_before, count_after)
}

func TestUserDelete(t *testing.T) {
    count_before := 0
    count_after := 0
    Users := [] models.User{}

    m := web.New()
    rooter(m)
    ts := httptest.NewServer(m)
    defer ts.Close()

    db.Find(&Users).Count(&count_before)
    req, _ := http.NewRequest("GET", ts.URL + "/user/delete/1", nil)
    req.Header.Set("Authorization", "Basic dXNlcjp1c2Vy")
    client := new(http.Client)
    client.Do(req)
    db.Find(&Users).Count(&count_after)

    assert.Equal(t, count_before - 1, count_after)
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
