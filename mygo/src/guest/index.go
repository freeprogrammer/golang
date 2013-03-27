package main

import (
	"guest/controller/content"
	"guest/controller/user"
	"github.com/dchest/captcha"
	"log"
	"net/http"
)

func main() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./view/CSS/"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./view/JS/"))))
	http.Handle("/image/", http.StripPrefix("/image", http.FileServer(http.Dir("./view/Images/"))))
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))

	http.HandleFunc("/reg", user.Reg)
	http.HandleFunc("/login", user.Login)
	http.HandleFunc("/guest", content.AddGuest)
	http.HandleFunc("/delguest", content.DelGuest)
	http.HandleFunc("/exit", user.Exit)
	http.HandleFunc("/", content.GetGuest)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe：", err)
	}
}
