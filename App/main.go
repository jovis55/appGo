package main

import (
	"fmt"
	"net/http"
	"REPASO_GO/App/functions"
	

)

func main(){

	// Ejemplo de slice de imágenes
	imageFiles := []string{"App/images/El Beso.jpg","App/images/El grito del Alma.jpg", "App/images/El Hijo del Hombre.jpg", "" }

	// Cantidad de imágenes a seleccionar aleatoriamente
	quantity := 2

	// Llamar a la función para obtener imágenes aleatorias sin repetición
	selectedImages := functions.getImages(imageFiles, quantity)

	// Imprimir las imágenes seleccionadas
	fmt.Println("Imágenes seleccionadas aleatoriamente:", selectedImages)
	//DEVOLVER UN  MENSAJE A LA PETICION DEL CLIENTE
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request){
		fmt.Fprintln(rw, "Hola Mundo")
	})
		
	fmt.Println("El servidor esta cotrriendo en puerto 3000")
	http.ListenAndServe("localhost:3000", nil)
}