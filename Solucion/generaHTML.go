package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	generaHTMLSol()
}

func generaHTMLSol() {
	fdIndexTemp, errIT := os.Open("indexTemplate.html")
	fdIndex, errIndex := os.OpenFile("sol/index2.html", os.O_CREATE|os.O_WRONLY, 0666)
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
		fmt.Println(linea)
		if linea == "				<h4 class='header-title'><b>Mejores Soluciónes en Algunas Iteraciones</b></h4>" {
			fmt.Fprintln(fdIndex, "Aqui va el contenido")
			fmt.Println("Aqui va el contenido")
		}
	}
}

func generaJSSol() {

}
