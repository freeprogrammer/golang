package main

import (
	"net/http"
	"fmt"
)

func testaa(w http.ResponseWriter, r *http.Request) {
	fmt.Println("aaa")
}
	
func main() {
	http.HandleFunc("/", testaa)
	http.ListenAndServe(":8080", nil)
	
}
