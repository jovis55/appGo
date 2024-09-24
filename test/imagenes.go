package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	imageDir := os.Args[1]
	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	var imageFiles []string
	validExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if validExtensions[ext] {
			imageFiles = append(imageFiles, filepath.Join(imageDir, file.Name()))
		}
	}

	// Barajar las imágenes y seleccionar una aleatoriamente
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(imageFiles), func(i, j int) {
		imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i]
	})

	if len(imageFiles) > 0 {
		selectedImage := imageFiles[0]
		fmt.Println("Imagen seleccionada:", filepath.Base(selectedImage))

		// Leer el contenido de la imagen seleccionada
		imgData, err := os.ReadFile(selectedImage)
		if err != nil {
			log.Fatal(err)
		}

		// Codificar la imagen a base64
		encoded := base64.StdEncoding.EncodeToString(imgData)
		fmt.Println("Imagen codificada en base64:", encoded)
	} else {
		fmt.Println("No se encontraron imágenes válidas en el directorio.")
	}
}
