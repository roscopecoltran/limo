package main

import (
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/jeffail/tunny"
)

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs+1) // numCPUs hot threads + one for async tasks.

	pool, _ := tunny.CreatePool(numCPUs, func(object interface{}) interface{} {
		input, _ := object.([]byte)

		// Do something that takes a lot of work
		output := input

		return output
	}).Open()

	defer pool.Close()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		input, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}

		// Send work to our pool
		result, _ := pool.SendWork(5000, input)

		w.Write(result.([]byte))
	})

	http.ListenAndServe(":8080", nil)
}