package main

import (
	"fmt"
	"os"
)

func main3() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Nombre del host desconocido"
	}
	fmt.Println(hostname)

}
