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

// Estructura para manejar la información de las imágeness
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
var puerto string

// Función para seleccionar imágenes aleatoriamente de una carpeta
func seleccionarImagenes(imageDir string, limit int) ([]ImageData, error) {
	files, err := os.ReadDir(imageDir)
	if err != nil {
		return nil, err
	}
	//SOLO SACARA LAS IMAGENES CON ESE FORMATO Y LAS PASA A MINUSCULAS
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

	//intercambia las posiciones del arreglo para sacar aleatoriamente, la 
	//la uncion que lo hace es shuffle, la funcion de intercambio es i j int
	// Esto intercambia los elementos en las posiciones i y j dentro de imageFiles. 
	//Es decir, el archivo en la posición i pasa a la posición j y viceversa.

	// Limitar el número de imágenes
	if len(imageFiles) > limit {
		imageFiles = imageFiles[:limit]
	}

	var images []ImageData

	//La función os.ReadFile se usa para leer
	// todo el contenido del archivo de imagen en forma de bytes (imgData).

	for _, imgPath := range imageFiles {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			log.Printf("Error al leer la imagen %s: %v", imgPath, err)
			continue
		}
		//El tipo MIME es necesario para indicar el formato de la imagen en la codificación Base64.
		ext := strings.ToLower(filepath.Ext(imgPath))
		mimeType := "image/png"
		if ext == ".jpg" || ext == ".jpeg" {
			mimeType = "image/jpeg"
		}
		//La imagen leída como bytes (imgData) se codifica en una cadena de caracteres
		// Base64 usando base64.StdEncoding.EncodeToString. Esto es útil 
		//para incrustar la imagen directamente en el HTML sin necesidad de tener un archivo separado.


		//e crea una URL en formato Data URI, que incluye el tipo MIME de la imagen y
		// los datos codificados en Base64
		encoded := base64.StdEncoding.EncodeToString(imgData)
		base64URL := template.URL("data:" + mimeType + ";base64," + encoded)
		images = append(images, ImageData{Base64: base64URL, Name: filepath.Base(imgPath)})

		//Se crea una nueva instancia de la estructura ImageData, asignando el valor 
		//codificado en Base64 (base64URL) y el nombre del archivo (filepath.Base(imgPath)).

	}

	return images, nil
}



//Esta función seleccionarNuevaPlantilla se utiliza para seleccionar 
//y alternar entre dos plantillas de forma condicional. La decisión de qué 
//plantilla devolver está basada en tres factores principales:

//La carpeta actual (carpeta).
//La cantidad actual de imágenes (currentCount).
//El estado anterior de carpeta y plantilla (previousFolder, previousCount, previousTemplate).

//CUando la caerpeta y las imagenes cambia hay una nueva plantilla
//y si cambia la plantilla entonces cambia la caerpeta y el numero de imagenes

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



//La función cargarPage es un handler de HTTP que carga una página web con una plantilla específica y datos personalizados, como las imágenes, el hostname,
// el tema y otros detalles de la materia. 

//El handler es una función o método que se encarga de manejar solicitudes HTTP 
//(como GET, POST, etc.) y responder a ellas.

// El handler debe construir y enviar una respuesta HTTP al cliente, que puede ser un archivo HTML

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
	if len(os.Args) != 4 {
		log.Fatal("Uso: go run servidor.go <carpeta> <plantilla>")
	}
	initialFolder = os.Args[1]
	initialTemplate = os.Args[2]
	puerto = os.Args[3]

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
	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
