/**
*	Autor: Martín Alejandro Pérez Güendulain
*	Correo: mperezguendulain@gmail.com
*	Descripción del Programa: Programa que trata de resolver el problema del Agente Viajero (TSP) ocupando Optimización ṕor Colonia de Hormigas
*	Ejecutar el programa de la siguiente forma:
*		go run hormigas.go < config.txt
*	Se puede generar el ejecutable de la siguiente forma:
*		go build hormigas.go
*	Y se ejecutaría así:
*		./hormigas < config.txt
*	Importante: Es necesario pasarle un arhivo con la configuración del problema
*	El archivo debe contener lo siguiente y ese orden:
*	Número de Generaciónes: 10
*	Número de Hormigas: 10000
*	K Vecinos más Cercanos: 5
*	Incremento: 0.05
*	Q: 4000
*	Número de Ciudades: 4
*	Ciudades:
*	500 0
*	462 0
*	44  0
*	176 0
*	Conexiónes:
*	0 1
*	0 2
*	0 3
**/

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

const (
	INFINITO = (1 << 32) - 1
)

type Config struct {
	NumCiudades     int
	NumHormigas     int
	NumGeneraciones int
	Knn             int
	Incremento      float64
	Q               int
	Ciudades        []Ciudad
	Conexiones      map[int][]int
}

type Ciudad struct {
	X int
	Y int
}

func (c Ciudad) String() string { return fmt.Sprintf("{x:%d, y:%d}", c.X, c.Y) }

type Ciudades []Ciudad

func (ciudades Ciudades) String() string {
	strCiudades := "["
	for i := 0; i < len(ciudades); i++ {
		if i != len(ciudades)-1 {
			strCiudades += fmt.Sprintf("%v, ", ciudades[i])
		} else {
			strCiudades += fmt.Sprintf("%v", ciudades[i])
		}
	}
	strCiudades += "]"
	return strCiudades
}

type Solucion struct {
	Camino []int
	Costo  float64
}

type Generacion []Solucion

func (g Generacion) Len() int           { return len(g) }
func (g Generacion) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g Generacion) Less(i, j int) bool { return g[i].Costo < g[j].Costo }

type infoIndexDistancias struct {
	Index     int
	Distancia float64
}

type ArrayInfoIndexDistancias []infoIndexDistancias

func (a ArrayInfoIndexDistancias) Len() int           { return len(a) }
func (a ArrayInfoIndexDistancias) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArrayInfoIndexDistancias) Less(i, j int) bool { return a[i].Distancia < a[j].Distancia }

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("Algoritmo de Colonia de Hormigas")
	var config = readConfig()
	printConfig(config)
	var generacion []Solucion
	mejorSol := Solucion{make([]int, config.NumCiudades), INFINITO}
	mejoresSoluciones := make([]Solucion, 0, config.NumGeneraciones)

	distancias, tau := getDistanciasYTau(config.Ciudades, config.Knn, config.Conexiones)

	for i := 0; i < config.NumGeneraciones; i++ {
		generacion = make([]Solucion, config.NumHormigas)

		for h := 0; h < config.NumHormigas; h++ {
			generacion[h].Camino = getNewCamino(distancias, tau, config.NumCiudades)
			generacion[h].Costo = getCostoCamino(generacion[h].Camino, distancias)
		}
		sort.Sort(Generacion(generacion))
		mejoresSoluciones = append(mejoresSoluciones, generacion[0])
		if generacion[0].Costo < mejorSol.Costo {
			mejorSol = generacion[0]
		}
		incrementaFeromonas(tau, config.Incremento)
		decrementaFeromonasEnMejoresCaminos(tau, generacion, config.Q)
	}
	if mejorSol.Costo == INFINITO {
		mejorSol = mejoresSoluciones[0]
	}
	printMejoresSoluciones(mejoresSoluciones)
	fmt.Println("Mejor Sol:", mejorSol)
	fmt.Println("--------------------------------")
	printMtz(distancias)
	fmt.Println("Costo Best Sol:", getCostoCamino(mejorSol.Camino, distancias))
	fmt.Println("--------------------------------")
	generaArchivosSol(mejoresSoluciones, mejorSol, config.Ciudades)
}

// Función que crea el archivo HTML con la solución del problema del agente viajero
func generaHTMLSol(mejoresSoluciones []Solucion, mejorSol Solucion, archivoHTML string) {
	fdIndexTemp, errIT := os.Open("Solucion/indexTemplate.html")
	fdIndex, errIndex := os.OpenFile(archivoHTML, os.O_CREATE|os.O_WRONLY, 0666)
	if errIT != nil {
		fmt.Println("Error:", errIT)
		return
	}
	defer fdIndexTemp.Close()
	if errIndex != nil {
		fmt.Println("Error:", errIndex)
		return
	}
	defer fdIndex.Close()

	scanner := bufio.NewScanner(fdIndexTemp)
	var linea string
	for scanner.Scan() {
		linea = scanner.Text()
		fmt.Fprintln(fdIndex, linea)
		if linea == "				<h4 class='header-title'><b>Mejores Soluciónes en Algunas Iteraciones</b></h4>" {
			fmt.Fprintln(fdIndex, "				<div class='steps_container' id='steps_container'>")
			for i := 0; i < len(mejoresSoluciones); i++ {
				if mejoresSoluciones[i].Costo == INFINITO {
					fmt.Fprintf(fdIndex, "					<div class='info-gen'>\n						<canvas id='canvas-gen-%d' width='500' height='500'></canvas>\n						<h4>Generación %d</h4>\n						<h5 class='camino-no-encontrado'>No se encontró un camino</h5>\n					</div>", i+1, i+1)
				} else {
					fmt.Fprintf(fdIndex, "					<div class='info-gen'>\n						<canvas id='canvas-gen-%d' width='500' height='500'></canvas>\n						<h4>Generación %d</h4>\n						<h5>Costo: %f</h5>\n					</div>", i+1, i+1, mejoresSoluciones[i].Costo)
				}
			}
			fmt.Fprintln(fdIndex, "\n				</div>")
			fmt.Fprintln(fdIndex, "				<div class='container-sol-final' id='container-sol-final'>")
			if mejorSol.Costo == INFINITO {
				fmt.Fprintf(fdIndex, "					<h4 class='header-title'><b>Mejor Solución</b></h4>\n						<div class='info-gen'>\n						<canvas id='canvas-best-sol' width='500' height='500'></canvas>\n						<h4>Mejor Solución</h4>\n						<h5 class='camino-no-encontrado'>No se encontró un camino</h5>\n					</div>\n")
			} else {
				fmt.Fprintf(fdIndex, "					<h4 class='header-title'><b>Mejor Solución</b></h4>\n						<div class='info-gen'>\n						<canvas id='canvas-best-sol' width='500' height='500'></canvas>\n						<h4>Mejor Solución</h4>\n						<h5>Costo: %f</h5>\n					</div>\n", mejorSol.Costo)
			}
			fmt.Fprintln(fdIndex, "				</div>")

		}
	}
}

// Función que crea el archivo Javascript con la solución del problema del agente viajero
func generaJSSol(mejoresSoluciones []Solucion, mejorSol Solucion, ciudades []Ciudad, archivoJs string) {
	fdJS, errJS := os.OpenFile(archivoJs, os.O_CREATE|os.O_WRONLY, 0666)
	if errJS != nil {
		fmt.Println("Error:", errJS)
		return
	}
	defer fdJS.Close()

	fmt.Fprintf(fdJS, "pasos = [\n")
	for i := 0; i < len(mejoresSoluciones); i++ {
		fmt.Fprintf(fdJS, "	%v,\n", Ciudades(getPasos(mejoresSoluciones[i].Camino, ciudades)))
	}
	fmt.Fprintf(fdJS, "];\n")

	fmt.Fprintf(fdJS, "mtz_avance_opt = [\n")
	for i := 0; i < len(mejoresSoluciones); i++ {
		fmt.Fprintf(fdJS, "	[%d, %f],\n", i+1, mejoresSoluciones[i].Costo)
	}
	fmt.Fprintf(fdJS, "];\n\n")

	fmt.Fprintf(fdJS, "drawPoints(%v, 'canvas-puntos-iniciales');\n", Ciudades(ciudades))
	fmt.Fprintf(fdJS, "drawSteps(pasos);\n")
	// no es ciudades, es mejorCamino
	fmt.Fprintf(fdJS, "drawBestSolution(%v);\n", Ciudades(getPasos(mejorSol.Camino, ciudades)))
}

// Función que dado un arreglo de posiciónes de Ciudades, retorna su respectivo arreglo de Ciudades
func getPasos(camino []int, ciudades []Ciudad) []Ciudad {
	pasos := make([]Ciudad, 0, len(camino))
	for i := 0; i < len(camino); i++ {
		pasos = append(pasos, ciudades[camino[i]])
	}
	return pasos
}

// Función que genera los archivos HTML y JS con la solución del problema del viajero y abre el navegador google-chrome para mostrar los resultados
func generaArchivosSol(mejoresSoluciones []Solucion, mejorSol Solucion, ciudades []Ciudad) {
	archivoHTML := "Solucion/index.html"
	archivoJS := "Solucion/js/init.js"
	os.Remove(archivoHTML)
	os.Remove(archivoJS)
	generaHTMLSol(mejoresSoluciones, mejorSol, archivoHTML)
	generaJSSol(mejoresSoluciones, mejorSol, ciudades, archivoJS)
	exec.Command("google-chrome", archivoHTML).Start()

}

// Función que imprime la ciudades pasadas
func printCiudades(ciudades []Ciudad) {
	for i := 0; i < len(ciudades); i++ {
		fmt.Println(i, "=>", ciudades[i])
	}
}

// Función que imprime las mejores soluciónes
func printMejoresSoluciones(mejoresSoluciones []Solucion) {
	for i := 0; i < len(mejoresSoluciones); i++ {
		fmt.Println(mejoresSoluciones[i])
	}
}

// Función que decremente las feromonas en los mejores caminos
func decrementaFeromonasEnMejoresCaminos(tau [][]float64, generacion []Solucion, Q int) {
	tamPorcion := len(generacion) / 4
	posFinal := len(tau) - 1
	for i := 0; i < tamPorcion; i++ {
		for j := 0; j < posFinal; j++ {
			tau[generacion[i].Camino[j]][generacion[i].Camino[j+1]] = tau[generacion[i].Camino[j]][generacion[i].Camino[j+1]] * (generacion[i].Costo / float64(Q))
		}
		tau[generacion[i].Camino[posFinal]][generacion[i].Camino[0]] = tau[generacion[i].Camino[posFinal]][generacion[i].Camino[0]] * (generacion[i].Costo / float64(Q))
	}
}

// Función que incrementa las feromonas en todos los caminos
func incrementaFeromonas(tau [][]float64, incremento float64) {
	numCiudades := len(tau)
	for i := 0; i < numCiudades; i++ {
		for j := 0; j < numCiudades; j++ {
			tau[i][j] *= (1 + incremento)
		}
	}
}

// Función que retorna un nuevo camino
func getNewCamino(distancias, tau [][]float64, numCiudades int) []int {
	visitados := make([]bool, numCiudades)
	camino := make([]int, 0, numCiudades)

	pos := rand.Intn(numCiudades)
	camino = append(camino, pos)
	visitados[pos] = true

	for i := 1; i < numCiudades; i++ {
		pos = getProxVisita(distancias[pos], tau[pos], visitados)
		visitados[pos] = true
		camino = append(camino, pos)
	}
	return camino
}

// Función que retorna la proxima ciudad a visitar
func getProxVisita(distanciasRow, tauRow []float64, visitados []bool) int {
	numCiudades := len(visitados)
	distanciasDifInf := make([]infoIndexDistancias, 0, numCiudades>>1)

	// Nos quedamos con las ciudades que sean diferentes de infinito
	for i := 0; i < numCiudades; i++ {
		if distanciasRow[i] != INFINITO {
			distanciasDifInf = append(distanciasDifInf, infoIndexDistancias{i, distanciasRow[i]})
		}
	}

	// Nos quedamos con las ciudades que aún no hayan sido visitadas
	distancias := make([]infoIndexDistancias, 0, len(distanciasDifInf))
	for i := 0; i < len(distanciasDifInf); i++ {
		if visitados[distanciasDifInf[i].Index] == false {
			distancias = append(distancias, distanciasDifInf[i])
		}
	}

	// Si hay buenos vecinos, hay que elegir uno de ellos con rank
	if len(distancias) > 0 {
		multiplicaDistanciasXTau(distancias, tauRow)
		sort.Sort(ArrayInfoIndexDistancias(distancias))
		iniRebanada := inicializaRank(len(distancias))
		return distancias[getPosition(rand.Float64(), iniRebanada, len(distancias))].Index

	} else { // Si no hay buenos vecinos, entonces que agarre aleatoriamente
		noVisitados := make([]int, 0, numCiudades)
		for i := 0; i < numCiudades; i++ {
			if visitados[i] == false {
				noVisitados = append(noVisitados, i)
			}
		}
		return noVisitados[rand.Intn(len(noVisitados))]
	}
}

// Función que se multiplica una fila de distancias por una fila de Tau, el resultado se queda en distancias
func multiplicaDistanciasXTau(distancias []infoIndexDistancias, tauRow []float64) {
	for i := 0; i < len(distancias); i++ {
		distancias[i].Distancia *= tauRow[distancias[i].Index]
	}
}

// Función para imprimir una matriz
func printMtz(mtz [][]float64) {
	for i := 0; i < len(mtz); i++ {
		for j := 0; j < len(mtz[0]); j++ {
			if mtz[i][j] == INFINITO {
				fmt.Printf("∞\t\t")
			} else {
				fmt.Printf("%.3f\t\t", mtz[i][j])
			}
		}
		fmt.Println()
	}
}

// Función para obtener la matriz de distancias y la matriz Tau
func getDistanciasYTau(ciudades []Ciudad, knn int, conexiones map[int][]int) ([][]float64, [][]float64) {
	distancias := getMtzDistancias(ciudades, conexiones)
	descartaMalosVecinos(distancias, knn)

	numCiudades := len(ciudades)
	for i := 0; i < numCiudades; i++ {
		for j := 0; j < numCiudades; j++ {
			if distancias[i][j] != INFINITO {
				distancias[j][i] = distancias[i][j]
			}
		}
	}

	buenosVecinos := getNumBuenosVecinos(distancias)
	tau := make([][]float64, numCiudades)
	for i := 0; i < numCiudades; i++ {
		tau[i] = make([]float64, numCiudades)
	}
	for i := 0; i < numCiudades; i++ {
		for j := 0; j < numCiudades; j++ {
			if distancias[i][j] != INFINITO {
				tau[i][j] = 1.0 / float64(buenosVecinos)
			}
		}
	}
	return distancias, tau
}

// Función que retorna el numero de elementos en la matriz de distancias que sean diferentes de INFINITO
func getNumBuenosVecinos(distancias [][]float64) int {
	numCiudades := len(distancias)
	cont := 0
	for i := 0; i < numCiudades; i++ {
		for j := 0; j < numCiudades; j++ {
			if distancias[i][j] != INFINITO {
				cont++
			}
		}
	}
	return cont
}

// Función para obtener la Matriz de Distancias
// Toda la matriz esta rellena de INFINITO, solo se cambian los valores dponde sí hay camino
func getMtzDistancias(ciudades []Ciudad, conexiones map[int][]int) [][]float64 {
	numCiudades := len(ciudades)
	distancias := make([][]float64, numCiudades)
	rowFillInf := make([]float64, numCiudades)
	for i := 0; i < numCiudades; i++ {
		rowFillInf[i] = INFINITO
	}
	for i := 0; i < numCiudades; i++ {
		distancias[i] = make([]float64, numCiudades)
		copy(distancias[i], rowFillInf)
	}

	for c, vecinos := range conexiones {
		for i := 0; i < len(vecinos); i++ {
			distancias[c][vecinos[i]] = getDistancia(ciudades[c], ciudades[vecinos[i]])
			distancias[vecinos[i]][c] = distancias[c][vecinos[i]]
		}
	}
	return distancias
}

// Se queda con los mejores knn vecinos de cada ciudad, los demás los pone INIFINITO
func descartaMalosVecinos(distancias [][]float64, knn int) {
	numCiudades := len(distancias)
	var stemp []float64
	for i := 0; i < numCiudades; i++ {
		stemp = make([]float64, numCiudades)
		copy(stemp, distancias[i])
		sort.Float64s(stemp)
		stemp = stemp[:1+knn]
		for j := 0; j < numCiudades; j++ {
			it := findInArray(distancias[i][j], stemp)
			if it == -1 {
				distancias[i][j] = INFINITO
			} else {
				stemp = append(stemp[:it], stemp[it+1:]...)
			}
		}
	}
}

// Funcion que encuentra un elemento en un slice. Retorna el indice donde se encontró, sino lo encontró retorna -1
func findInArray(num float64, numeros []float64) int {
	tamArray := len(numeros)
	for i := 0; i < tamArray; i++ {
		if num == numeros[i] {
			return i
		}
	}
	return -1
}

// Función para obtener el costo total de un camino
func getCostoCamino(camino []int, mtzDistancias [][]float64) float64 {
	pos_final := len(camino) - 1
	costo_final := mtzDistancias[camino[0]][camino[pos_final]]
	var distTemp float64
	if costo_final == INFINITO {
		return INFINITO
	}
	for pos := 0; pos < pos_final; pos++ {
		distTemp = mtzDistancias[camino[pos]][camino[pos+1]]
		if distTemp == INFINITO {
			return INFINITO
		}
		costo_final += distTemp
	}
	return costo_final
}

// Función para obtener la distancia entre 2 ciudades
func getDistancia(a, b Ciudad) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2.0) + math.Pow(float64(a.Y-b.Y), 2.0))
}

// Función para imprimir la configuración inicial
func printConfig(config Config) {
	fmt.Println("Numero de Ciudades:", config.NumCiudades)
	fmt.Println("Numero de Generaciónes:", config.NumGeneraciones)
	fmt.Println("Numero de Hormigas:", config.NumHormigas)
	fmt.Println("K Vecinos más Cercanos:", config.Knn)
	fmt.Println("Incremento:", config.Incremento)
	fmt.Println("Q:", config.Q)
	fmt.Println("Ciudades:", config.Ciudades)
	fmt.Println("Conexiones:", config.Conexiones)
}

// Función para leer la configuración inicial
func readConfig() Config {
	var config = Config{}
	fmt.Scanf("Número de Generaciónes: %d", &config.NumGeneraciones)
	fmt.Scanf("Número de Hormigas: %d", &config.NumHormigas)
	fmt.Scanf("K Vecinos más Cercanos: %d", &config.Knn)
	fmt.Scanf("Incremento: %f", &config.Incremento)
	fmt.Scanf("Q: %d", &config.Q)
	fmt.Scanf("Número de Ciudades: %d\nCiudades:\n", &config.NumCiudades)

	config.Ciudades = make([]Ciudad, 0, config.NumCiudades)
	var x, y int
	config.Conexiones = make(map[int][]int, 0)
	for i := 0; i < config.NumCiudades; i++ {
		fmt.Scanf("%d %d\n", &x, &y)
		config.Ciudades = append(config.Ciudades, Ciudad{x, y})
		config.Conexiones[i] = make([]int, 0)
	}
	fmt.Scanf("Conexiónes:\n")
	var ciudadA, ciudadB int
	for {
		n, _ := fmt.Scanf("%d %d", &ciudadA, &ciudadB)
		if n == 0 {
			break
		}
		config.Conexiones[ciudadA] = append(config.Conexiones[ciudadA], ciudadB)
	}

	return config
}

/**
* Función auxiliar de Rank que retorna la posición donde cayó la bola (t), en la ruleta (ini_rebanada)
* t, numero entre 0 y 1
* ini_rebanada es un arreglo donde en cada celda tiene el numero inicial de una rebanada
* tam_gen es el tamaño del arreglo ini_rebada
**/
func getPosition(t float64, ini_rebanada []float64, tam_gen int) int {
	for i := 0; i < tam_gen; i++ {
		if ini_rebanada[i] > t {
			return i - 1
		}
	}
	return tam_gen - 1
}

// Función que inicializa el Rank
func inicializaRank(tam_gen int) []float64 {
	rand.Seed(time.Now().Unix())
	sumatoria_rank := sumatoria(tam_gen)
	ini_rebanada := make([]float64, 0, tam_gen)
	ini := 0.0

	for i := 0; i < tam_gen; i++ {
		porcion := float64(tam_gen-i) / float64(sumatoria_rank)
		ini_rebanada = append(ini_rebanada, ini)
		ini += porcion
	}
	return ini_rebanada
}

// Función que calcula la sumatoria de N
// Ej: sumatoria(N) = N+N-1+N-2+N-3+...
func sumatoria(n int) int {
	sumatoria := 0
	for n > 0 {
		sumatoria += n
		n--
	}
	return sumatoria
}
