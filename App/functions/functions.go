package function2

import(
	"math/rand"
	"time"

)


func getImages(images []string, cantidad int) []string {
	// Asegurarse de que la cantidad solicitada no sea mayor que la cantidad de imágenes disponibles
	if cantidad > len(images) {
		cantidad = len(images)
	}

	// Mezclar el slice de imágenes de forma aleatoria
	// Usa el tiempo actual en nanosegundos como semilla, lo que asegura 
	//que los números generados (y por lo tanto las imágenes seleccionadas) 
	//sean diferentes cada vez que ejecutas el programa.
	//LA SEMILLA DEBE DE CAMBIAR CADA VEZ

	rand.Seed(time.Now().UnixNano()) // Semilla para la aleatoriedad
	rand.Shuffle(len(images), func(i, j int) {
		images[i], images[j] = images[j], images[i]
	})

	// Devolver las primeras `quantity` imágenes sin repetir
	return images[:cantidad]
}