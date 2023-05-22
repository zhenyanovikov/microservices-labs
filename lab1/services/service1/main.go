package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Person struct {
	ID   int
	Name string
}

func main() {
	// run http server with 4 CRUD endpoints
	httpMux := http.NewServeMux()

	prefix := "/service1"

	data := map[int]Person{
		1: {
			ID:   1,
			Name: "John",
		},
		2: {
			ID:   2,
			Name: "Jane",
		},
		3: {
			ID:   3,
			Name: "Joe",
		},
	}

	httpMux.HandleFunc(prefix+"/create", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")

		data[len(data)+1] = Person{
			ID:   len(data),
			Name: name,
		}

		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/read", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil && idStr != "" {
			return
		}

		var result string

		if id != 0 {
			result = data[id].Name
		} else {
			for key, person := range data {
				result += fmt.Sprintf("%d) %s\n", key, person.Name)
			}
		}

		w.Write([]byte(result))
	})

	httpMux.HandleFunc(prefix+"/update", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}

		data[id] = Person{
			ID:   id,
			Name: name,
		}

		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/delete", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}

		delete(data, id)

		w.Write([]byte("OK"))
	})

	panic(http.ListenAndServe(":8080", httpMux))
}
