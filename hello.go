/* package main

import (
	"fmt"
	"net/http"
	"database/sql"


	_ "github.com/go-sql-driver/mysql"
)

func main() {

fmt.Println("Go mysql tuto")

	// connexion db
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/zinzinebis")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully conected tu mysql database")


    fmt.Println("Building REST API in Go")
	// routes
	mux:= http.NewServeMux()

	mux.HandleFunc("GET /emission", func(w http.ResponseWriter, r *http.Request) {
		result, err := db.Query("SELECT * FROM emission")

		if err != nil {
			panic(err.Error())
		}
		http.Handle(result)


	})

	mux.HandleFunc("POST /emission", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Post a new emission")

	})

	mux.HandleFunc("GET /emission/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintln(w, "Return a single emission whith id : ", id)

	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}
} */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)
type Emission struct {
  ID string `json:"id"`
  Titre string `json:"titre"`
  Datepub string `json:"datepub"`
  Duree string `json:"duree"`
  Url string `json:"url"`
  Descriptif string `json:"descriptif"`
}
var db *sql.DB
var err error
func main() {
  db, err = sql.Open("mysql", "root:@/zinzinebis")
  
  if err != nil {
    panic(err.Error())
  }
defer db.Close()
router := mux.NewRouter()
router.HandleFunc("/emissions", getEmissions).Methods("GET")
  router.HandleFunc("/emissions", createEmission).Methods("POST")
  router.HandleFunc("/emissions/{id}", getEmission).Methods("GET")
  router.HandleFunc("/emissions/{id}", updateEmission).Methods("PUT")
  router.HandleFunc("/emissions/{id}", deleteEmission).Methods("DELETE")
http.ListenAndServe(":8000", router)
}

func getEmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var emissions []Emission
	result, err := db.Query("SELECT id, titre, datepub, duree, url, descriptif from emission")
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
	  var emission Emission
	  err := result.Scan(&emission.ID, &emission.Titre, &emission.Datepub, &emission.Duree, &emission.Url, &emission.Descriptif)
	  if err != nil {
		panic(err.Error())
	  }
	  emissions = append(emissions, emission)
	}
	json.NewEncoder(w).Encode(emissions)
  }
  func createEmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO emission(titre) VALUES(?)")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	titre := keyVal["titre"]
	_, err = stmt.Exec(titre)
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "New emission was created")
  }
  func getEmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, titre FROM emission WHERE id = ?", params["id"])
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var emission Emission
	for result.Next() {
	  err := result.Scan(&emission.ID, &emission.Titre, &emission.Datepub, &emission.Duree, &emission.Url, &emission.Descriptif)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(emission)
  }
  func updateEmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE emission SET titre = ? WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitre := keyVal["titre"]
	_, err = stmt.Exec(newTitre, params["id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Emission with ID = %s was updated", params["id"])
  }
  func deleteEmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM emission WHERE id = ?")
	if err != nil {
	  panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Emission with ID = %s was deleted", params["id"])
  }

