package main

import (
  "os"
  "time"
  "io/ioutil"
  "net/http"
  "log"
  "github.com/gorilla/mux"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func CreateDatabase(db_path string) {
  dat := []byte("ENVIRONMENT,USERNAME,VERSION,TIME\n")
  err := ioutil.WriteFile(db_path, dat, 0644)
  check(err)
}

func ReadDatabase(w http.ResponseWriter, r *http.Request) {
  db_path := "db.csv"
  dat, err := ioutil.ReadFile(db_path)
  check(err)
  w.Write([]byte(string(dat)))
}

func WriteToDatabase(w http.ResponseWriter, r *http.Request) {
  db_path := "db.csv"
  vars := mux.Vars(r)
  log.Println(vars["env"] + "\t" + vars["user"] + "\t" + vars["version"])
  f, err := os.OpenFile(db_path, os.O_APPEND|os.O_WRONLY, 0644)
  check(err)
  defer f.Close()
  time_now := time.Now().Format(time.RFC3339)
  if _, err = f.WriteString(vars["env"] + "," + vars["user"] + "," + vars["version"] + "," + time_now + "\n"); err != nil {
    panic(err)
  }
}

func main() {
  db_path := "db.csv"
  if _, err := os.Stat(db_path); os.IsNotExist(err) {
    CreateDatabase(db_path)
  }
  r := mux.NewRouter()
  r.HandleFunc("/{env}/{user}/{version}", WriteToDatabase).Methods("POST")
  r.HandleFunc("/", ReadDatabase).Methods("GET")
  log.Fatal(http.ListenAndServe(":8000", r))
}
