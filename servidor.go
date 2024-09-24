package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageData struct {
	Base64 template.URL
	Name   string
}

type PageData struct {
	Hostname     string
	Images       []ImageData
	Subject      string
	Materia      string
	Participants []string
	Date         string
}

func cargarPage(images []ImageData, hostname string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Hostname:     hostname,
			Images:       images,
			Subject:      "Camionetas y SUVs",
			Materia:      "Computación en la nube",
			Participants: []string{"Johana Palacio", "Alejandro Zapata"},
			Date:         "2024-2",
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error al ejecutar la plantilla: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Se debe proporcionar la ruta de la carpeta de imágenes como argumento.")
	}

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

	// Barajar las imágenes para seleccionar aleatoriamente
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(imageFiles), func(i, j int) {
		imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i] // Intercambiar posiciones
	})

	if len(imageFiles) > 4 {
		imageFiles = imageFiles[:4]
	}

	var images []ImageData
	for _, imgPath := range imageFiles {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			log.Printf("Error al leer la imagen %s: %v", imgPath, err)
			continue
		}

		ext := strings.ToLower(filepath.Ext(imgPath))
		mimeType := "image/png"
		if ext == ".jpg" || ext == ".jpeg" {
			mimeType = "image/jpeg"
		}

		encoded := base64.StdEncoding.EncodeToString(imgData)
		imageName := filepath.Base(imgPath)
		base64URL := template.URL("data:" + mimeType + ";base64," + encoded)
		images = append(images, ImageData{Base64: base64URL, Name: imageName})
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Nombre del host desconocido"
	}

	http.HandleFunc("/", cargarPage(images, hostname))
	fmt.Println("Servidor iniciado en el puerto 9090")

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
