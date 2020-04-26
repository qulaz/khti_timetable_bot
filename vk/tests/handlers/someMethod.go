package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

// Возвращает строку key=value; из параметров пост запроса
func SomeMethodHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// Sorting form keys
	var keys []string
	for k := range r.Form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	responseString := ""
	for _, key := range keys {
		responseString += fmt.Sprintf("%s=%q;", key, r.Form.Get(key))
	}

	if _, err := fmt.Fprintln(w, responseString); err != nil {
		log.Fatal(err)
	}
}
