package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"stixi/back/hlp"
)



type Connect struct {
	Key     string
	Crt     string
	Host     string
	Port     string
}


func main() {
	var connect Connect
	// Чтение переменных из файла настроек ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	i := hlp.InitS{File: "config.ini", Sect: "Connect"}
	if err := i.InitF(&connect); err != nil {
		hlp.Add2Log("conn.go", "fail to read ini: "+fmt.Sprint(err))
	}

	var addr = flag.String("addr", connect.Host+":"+connect.Port, "http service address")

	flag.Parse()
	go h.run()
	//http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)

	var httpErr error

	if _, err := os.Stat(connect.Crt); err == nil {
		fmt.Println("file ", "crt found switching to https")
		if httpErr = http.ListenAndServeTLS(*addr, connect.Crt, connect.Key, nil); httpErr != nil {
			log.Fatal("The process exited with https error: ", httpErr.Error())
		}
	} else {
		httpErr = http.ListenAndServe(*addr, nil)
		fmt.Println("file", "crt not found")
		if httpErr != nil {
			log.Fatal("The process exited with http error: ", httpErr.Error())
		}
	}

}
