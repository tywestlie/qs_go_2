package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "encoding/json"


    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(db string) {
  connectionString :=
      fmt.Sprintf("%s", db)

  var err error
  a.DB, err = sql.Open("postgres", connectionString)
  if err != nil {
      log.Fatal(err)
  }

  a.Router = mux.NewRouter()
  a.initializeRoutes()
}

func (a *App) Run(addr string) {
  log.Fatal(http.ListenAndServe(":3000", a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) initializeRoutes() {
  a.Router.HandleFunc("/food/{id:[0-9]+}", a.getFood).Methods("GET")
}

func (a *App) getFood(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid food ID")
        return
    }

    p := food{ID: id}
    if err := p.getFood(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Food not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}
