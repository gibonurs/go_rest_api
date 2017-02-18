package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/user"
)

type App struct {
  Router *mux.Router
  Mongo *MongoConnection
}

type MongoConnection struct {
  Session *mgo.Session
}

func(a* App) GetMongoSession() *mgo.Session {
  return a.Mongo.Session.Copy()
}

func(a *App) Initialize() {
  fmt.Println("Initialize")
  session, err := mgo.Dial("127.0.0.1:27017")
  if err != nil {
    panic(err)
  }

  session.SetMode(mgo.Monotonic, true)

  a.Mongo = &MongoConnection{session}
  a.Router = mux.NewRouter()
  
  a.SetupRoutes()
}

func(a *App) SetupRoutes() {
  fmt.Println("set up routes")

  a.Router.Handle("/user", user.NewUserRouter(a.getSubrouter("/user")))
  
  http.Handle("/", a.Router)
}

func(a *App) getSubrouter(path string) (*mgo.Session, *mux.Router) {
  return a.GetMongoSession(), a.Router.PathPrefix(path).Subrouter()
}

func(a *App) Run() {
  fmt.Println("Run")
  defer a.Mongo.Session.Close()

  log.Println("Listening on port 1337...")
  if err := http.ListenAndServe(":1337", nil); err != nil {
      log.Fatal("http.ListenAndServe: ", err)
  }
}