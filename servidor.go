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

// Estructura para manejar la información de las imágenes
type ImageData struct {
	Base64 template.URL
	Name   string
}

// Estructura para la página con imágenes, tema dinámico y materia
type PageData struct {
	Hostname     string
	Images       []ImageData
	Subject      string // El tema (Camionetas y SUVs o Pinturas emblemáticas)
	Materia      string // Nueva propiedad que representa la materia
	Date         string
	Participants []string
}

var previousTemplate string // Para evitar que la misma plantilla se repita consecutivamente
var previousFolder string   // Para evitar que la misma carpeta y número de imágenes se repita
var previousCount int       // Para almacenar el conteo de imágenes mostradas
var isFirstRun = true       // Bandera para la primera ejecución
var initialFolder string    // Carpeta pasada en la ejecución
var initialTemplate string  // Plantilla pasada en la ejecución

// Función para seleccionar imágenes aleatoriamente de una carpeta
func seleccionarImagenes(imageDir string, limit int) ([]ImageData, error) {
	files, err := os.ReadDir(imageDir)
	if err != nil {
		return nil, err
	}

	var imageFiles []string
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if validExtensions[ext] {
			imageFiles = append(imageFiles, filepath.Join(imageDir, file.Name()))
		}
	}

	// Barajar las imágenes aleatoriamente
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(imageFiles), func(i, j int) {
		imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i]
	})

	// Limitar el número de imágenes
	if len(imageFiles) > limit {
		imageFiles = imageFiles[:limit]
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
		base64URL := template.URL("data:" + mimeType + ";base64," + encoded)
		images = append(images, ImageData{Base64: base64URL, Name: filepath.Base(imgPath)})
	}

	return images, nil
}

// Función para alternar entre las plantillas
func seleccionarNuevaPlantilla(carpeta string, currentCount int) string {
	if previousFolder == carpeta && previousCount == currentCount {
		// Cambiar plantilla si es la misma carpeta y el mismo conteo de imágenes
		if previousTemplate == "1" {
			previousTemplate = "2"
			previousCount = 3
			return "2"
		} else {
			previousTemplate = "1"
			previousCount = 4
			return "1"
		}
	} else if previousFolder != carpeta && previousCount != currentCount {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumberImagesPlantilla := r.Intn(2) + 3
		if randomNumberImagesPlantilla == 3 {
			previousCount = 3
			return "2"
		} else {
			previousCount = 4
			return "1"
		}
	} else if previousFolder != carpeta && previousCount == currentCount {
		if currentCount == 3 {
			previousCount = 3
			return "2"
		} else {
			previousCount = 4
			return "1"
		}

	} else {
		if currentCount == 3 {
			previousCount = 3
			return "2"
		} else {
			previousCount = 4
			return "1"
		}
	}
}

// Función para cargar la página con la plantilla elegida, el tema asociado y la materia
func cargarPage(images []ImageData, hostname string, plantillaElegida string, tema string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Hostname:     hostname,
			Images:       images,
			Subject:      tema,
			Materia:      "Computación en la nube", // Añadido el campo Materia
			Date:         "2024-2",
			Participants: []string{"Johana Palacio", "Alejandro Zapata"},
		}
		tmpl := template.Must(template.ParseFiles(plantillaElegida))
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error al ejecutar la plantilla: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	}
}

func main() {
	// Obtener los parámetros iniciales desde la línea de comandos
	if len(os.Args) != 3 {
		log.Fatal("Uso: go run servidor.go <carpeta> <plantilla>")
	}
	initialFolder = os.Args[1]
	initialTemplate = os.Args[2]

	// Inicializar previousTemplate y previousFolder con los valores pasados por línea de comandos
	previousTemplate = initialTemplate
	previousFolder = initialFolder

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var carpeta, plantilla, tema string
		var imageLimit int

		var plantillaElegida string

		// Primera ejecución con parámetros
		if isFirstRun {
			carpeta = initialFolder
			plantilla = initialTemplate
			previousTemplate = plantilla

			// Asignar el tema según la carpeta
			if carpeta == "camionetas" {
				tema = "Camionetas y SUVs"
			} else if carpeta == "arte" {
				tema = "Pinturas emblemáticas"
			}

			isFirstRun = false // La primera ejecución ya ocurrió
		} else {
			// Recargas aleatorias
			selecciones := []string{"camionetas", "arte"}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			carpeta = selecciones[r.Intn(len(selecciones))]

			// Asignar el tema según la nueva carpeta
			if carpeta == "camionetas" {
				tema = "Camionetas y SUVs"
			} else if carpeta == "arte" {
				tema = "Pinturas emblemáticas"
			}
			// Generar un número aleatorio entre 3 y 4
			randomNumberImagesPlantilla := rand.Intn(2) + 3

			plantilla = seleccionarNuevaPlantilla(carpeta, randomNumberImagesPlantilla)

		}

		if plantilla == "1" {
			plantillaElegida = "templates/index.html"
			previousCount = 4
			imageLimit = 4
		} else if plantilla == "2" {
			plantillaElegida = "templates/index2.html"
			previousCount = 3
			imageLimit = 3
		}

		// Cargar las imágenes de la carpeta seleccionada
		images, err := seleccionarImagenes(carpeta, imageLimit)
		if err != nil {
			http.Error(w, "Error al cargar imágenes", http.StatusInternalServerError)
			return
		}

		// Obtener el nombre del host
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "Nombre del host desconocido"
		}

		// Actualizar variables de seguimiento después de mostrar la página
		previousFolder = carpeta
		previousTemplate = plantilla // Actualizar la plantilla usada

		// Cargar la página con la nueva selección
		cargarPage(images, hostname, plantillaElegida, tema)(w, r)

	})

	fmt.Println("Servidor iniciado en el puerto 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
