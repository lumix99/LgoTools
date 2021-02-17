package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lumix99/LgoTools/example/bakuteh/controller"

	"github.com/lumix99/LgoTools/bakuteh"
)

func main() {
	//var addr = flag.String("addr", ":8080", "http service address")
	defer fmt.Println("process ready to exit")
	bk := bakuteh.New()
	bk.HandleServer("/", &controller.Chat{})

	http.HandleFunc("/", bk.ServerWs)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
