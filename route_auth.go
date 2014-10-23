package main

import (
  "fmt"
  "net/http"
  "time"
  "./data"
)

// GET /
// index page
func index(writer http.ResponseWriter, request *http.Request) {
  loggedin(writer, request)
  t := parseTemplateFiles("layout", "public.navbar", "index")
  t.Execute(writer, "helloz")
}

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
  t := parseTemplateFiles("login.layout", "public.navbar", "login")
  t.Execute(writer, nil)
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
  t := parseTemplateFiles("login.layout", "public.navbar", "signup")
  t.Execute(writer, nil)  
}

// POST /signup
// Create the user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
  err := request.ParseForm()
  if err != nil {
    fmt.Println("err", err)
  }
  user := data.User{
    Name: request.PostFormValue("name"),
    Email: request.PostFormValue("email"),
    Password: request.PostFormValue("password"),    
  }
  if err := user.Create(); err != nil {
    fmt.Println(err, "Cannot create user.")
  }
  http.Redirect(writer, request, "/login", 302)
}


// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {  
  err := request.ParseForm()
  user, err := data.UserByEmail(request.PostFormValue("email"))
  if err != nil {
    fmt.Println("Cannot find user")
  }

  if user.Password == data.Encrypt(request.PostFormValue("password")) {
    session, err := user.CreateSession()
    if err != nil {
      fmt.Println("Cannot create session")
    }
    cookie := http.Cookie{
      Name:      "_cookie", 
      Value:     session.Uuid,
      Expires:   time.Now().Add(90 * time.Minute),
      HttpOnly:  true,
    }
    http.SetCookie(writer, &cookie)
    http.Redirect(writer, request, "/", 302)
  } else {
    http.Redirect(writer, request, "/login", 302)
  }
  
}

func logout(writer http.ResponseWriter, request *http.Request) {
  cookie, err := request.Cookie("_cookie")
  if err != http.ErrNoCookie {
    fmt.Println(err, "Failed to get cookie")
    session := data.Session{Uuid: cookie.Value}
    session.DeleteByUUID()    
  }  
  http.Redirect(writer, request, "/", 302)
}


func loggedin(writer http.ResponseWriter, request *http.Request)(authenticated bool){
  cookie, err := request.Cookie("_cookie")
  if err == http.ErrNoCookie {
    http.Redirect(writer, request, "/login", 302)
  } else {
    session := data.Session{Uuid: cookie.Value}
    if ok, _ := session.Check(); !ok {
      http.Redirect(writer, request, "/login", 302)
    }
  }  
  return
}

func index2(writer http.ResponseWriter, request *http.Request) {
  fmt.Println("request", request)
  fmt.Println("url", request.URL)
  fmt.Println("url.path", request.URL.Path)
  fmt.Println("url.rawquery", request.URL.RawQuery)
  fmt.Println("url.fragment", request.URL)
  vals := request.URL.Query()
  fmt.Println("url.query", vals.Get("id"))
  t := parseTemplateFiles("layout", "public.navbar", "index")
  t.Execute(writer, "helloz")
}