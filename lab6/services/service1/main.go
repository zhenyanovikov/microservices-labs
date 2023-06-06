package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const prefix = "/service1"
const service2URL = "http://service2:80"
const version = "0.4"

type Person struct {
	bun.BaseModel `bun:"table:persons"`
	ID            int `bun:",pl,nullzero"`
	Name          string
}

func main() {
	httpMux := http.NewServeMux()

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("POSTGRES_URL"))))
	db := bun.NewDB(sqlDB, pgdialect.New())

	_, err := db.NewSelect().ColumnExpr("1").Exec(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Version: " + version)
	fmt.Println("Connected to database")

	isBroken := false

	// add 404 handler

	httpMux.HandleFunc(prefix+"/chain-test", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		//	make a request to service2
		req, err := http.NewRequest("GET", service2URL+"/service2/read", nil)
		if err != nil {
			fmt.Println("Error creating request")
			return
		}

		if _, err = io.Copy(w, req.Body); err != nil {
			fmt.Println("Error copying body")
			return
		}
	})

	httpMux.HandleFunc(prefix+"/broke", func(w http.ResponseWriter, r *http.Request) {
		isBroken = true
		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/revive", func(w http.ResponseWriter, r *http.Request) {
		isBroken = false
		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/health", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/create", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		name := r.URL.Query().Get("name")

		person := Person{
			Name: name,
		}

		db.NewInsert().Model(&person).Exec(r.Context())

		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/read", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil && idStr != "" {
			return
		}

		var data []Person

		q := db.NewSelect().
			Model(&data)

		if id != 0 {
			q.Where("id = ?", id)
		}

		if err = q.Scan(r.Context()); err != nil {
			panic(err)
		}

		var result string

		for i := range data {
			result += fmt.Sprintf("%d) %s\n", data[i].ID, data[i].Name)
		}

		w.Write([]byte(result))
	})

	httpMux.HandleFunc(prefix+"/update", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		idStr := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}

		db.NewUpdate().
			Model((*Person)(nil)).
			Set("name = ?", name).
			Where("id = ?", id).
			Exec(r.Context())

		w.Write([]byte("OK"))
	})

	httpMux.HandleFunc(prefix+"/delete", func(w http.ResponseWriter, r *http.Request) {
		checkBroken(w, isBroken)

		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			panic(err)
		}

		db.NewDelete().
			Model((*Person)(nil)).
			Where("id = ?", id).
			Exec(r.Context())

		w.Write([]byte("OK"))
	})

	panic(http.ListenAndServe(":8080", httpMux))
}

func checkBroken(w http.ResponseWriter, isBroken bool) {
	if isBroken {
		<-time.After(10 * time.Second)

		// return 500
		w.WriteHeader(http.StatusInternalServerError)
	}
}
