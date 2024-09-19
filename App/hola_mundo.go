package main

import (
	"fmt"
	"net/http"
)

func main(){
	//DEVOLVER UN  MENSAJE A LA PETICION DEL CLIENTE
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request){
		fmt.Fprintln(rw, "Hola Mundo")
	})
		
	fmt.Println("El servidor esta cotrriendo en puerto 3000")
	http.ListenAndServe("localhost:3000", nil)
}