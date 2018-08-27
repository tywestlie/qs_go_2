package main_test

import (
    "os"
    "testing"
    "log"
    "net/http"
    "net/http/httptest"
    "encoding/json"

    "test"
)

var a main.App

func TestMain(m *testing.M) {
    a = main.App{}
    a.Initialize("dbname=test_qs_go sslmode=disable")

    ensureTableExists()

    code := m.Run()

    clearTable()

    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM foods")
    a.DB.Exec("ALTER SEQUENCE foodss_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/foods", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestGetNonExistentFood(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/foods/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Food not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Food not found'. Got '%s'", m["error"])
    }
}



const tableCreationQuery = `CREATE TABLE IF NOT EXISTS foods
(
id SERIAL,
name TEXT NOT NULL,
calories NUMERIC(10,2) NOT NULL DEFAULT 0.00,
CONSTRAINT foods_pkey PRIMARY KEY (id)
)`
