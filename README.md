# Problema del Agente Viajero solucionado con Optimización por Colonia de Hormigas

Programa que proporciona una solución subóptima al problema del Agente Viajero (TSP) ocupando Optimización por Colonia de Hormigas.

![Optimización por Colonias de Hormigas](/ant-colony-go.png)

![Optimización por Colonias de Hormigas](/res-ant-colony1.png)

![Optimización por Colonias de Hormigas](/res-ant-colony2.png)

## Optimización por Colonia de Hormigas
**Optimización por Colonia de Hormigas** es una tecnica probabilistica para optimización. La fuente inspiradora de la optimización por colinias de hormigas es el comportamiento forrajeo (buscar ampliamente alimentos o provisiones) de colonias de hormigas reales. Buscando el camino óptimo en el gráfico basado en el comportamiento de las hormigas buscando un camino entre su colonia y la fuente de alimento.

## Problema del Agente viajero (TSP)
En el Problema del Agente Viajero - TSP (Travelling Salesman Problem), el objetivo es visitar todas las ciudades de un conjunto de ciudades, pasando por cada ciudad solamente una vez, volviendo al punto de partida, y que además minimice el costo total de la ruta.

### Requisitos
 - Google Chrome con soporte canvas

### Compilación
	go build hormigas.go

### Ejecución
	./hormigas < config.txt

#### **Importante**
Es necesario pasarle un archivo con la configuración del problema.  El archivo debe contener lo siguiente y ese orden:

    Número de Generaciónes: 10
    Número de Hormigas: 10000
    K Vecinos más Cercanos: 5
    Incremento: 0.05
    Q: 4000
    Número de Ciudades: 4
    Ciudades:
    500 0
    462 0
    44  0
    176 0
    Conexiónes:
    0 1
    0 2
    0 3
