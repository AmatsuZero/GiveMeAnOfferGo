package Graph

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
	"strings"
)

type AdjacencyMatrix struct {
	vertices []Vertex
	weights  [][]float64
}

func NewAdjacencyMatrix() *AdjacencyMatrix {
	return &AdjacencyMatrix{
		vertices: make([]Vertex, 0),
		weights:  make([][]float64, 0),
	}
}

func (matrix *AdjacencyMatrix) CreateVertex(data Objects.EquatableObject) Vertex {
	vertex := Vertex{
		Index: len(matrix.vertices),
		Data:  data,
	}
	matrix.vertices = append(matrix.vertices, vertex)
	for i := 0; i < len(matrix.weights); i++ {
		matrix.weights[i] = append(matrix.weights[i], Objects.InvalidResult)
	}
	row := make([]float64, len(matrix.vertices))
	for i := 0; i < len(matrix.vertices); i++ {
		row[i] = Objects.InvalidResult
	}
	matrix.weights = append(matrix.weights, row)
	return vertex
}

func (matrix *AdjacencyMatrix) AddDirectedEdge(source Vertex, destination Vertex, weight float64) {
	matrix.weights[source.Index][destination.Index] = weight
}

func (matrix *AdjacencyMatrix) AddUndirectedEdge(source Vertex, destination Vertex, weight float64) {
	matrix.AddDirectedEdge(source, destination, weight)
	matrix.AddDirectedEdge(destination, source, weight)
}

func (matrix *AdjacencyMatrix) Add(edgeType EdgeType, source Vertex, destination Vertex, weight float64) {
	switch edgeType {
	case EdgeTypeDirected:
		matrix.AddDirectedEdge(source, destination, weight)
	case EdgeTypeUndirected:
		matrix.AddUndirectedEdge(source, destination, weight)
	}
}

func (matrix *AdjacencyMatrix) Edges(source Vertex) []Edge {
	edges := make([]Edge, 0)
	for column := 0; column < len(matrix.weights); column++ {
		weight := matrix.weights[source.Index][column]
		if weight != Objects.InvalidResult {
			edges = append(edges, Edge{
				Source:      source,
				Destination: matrix.vertices[column],
				Weight:      weight,
			})
		}
	}
	return edges
}

func (matrix *AdjacencyMatrix) Weight(source Vertex, destination Vertex) float64 {
	return matrix.weights[source.Index][destination.Index]
}

func (matrix *AdjacencyMatrix) String() string {
	verticesDesc := ""
	for _, vertex := range matrix.vertices {
		verticesDesc += fmt.Sprintf("%v\n", vertex)
	}
	grid := make([]string, 0)
	for i := 0; i < len(matrix.weights); i++ {
		row := ""
		for j := 0; j < len(matrix.weights); j++ {
			value := matrix.weights[i][j]
			if value == Objects.InvalidResult {
				row += "ø\t\t"
			} else {
				row += fmt.Sprintf("%v\t", value)
			}
		}
		grid = append(grid, row)
	}
	edgeDesc := strings.Join(grid, "\n")
	return fmt.Sprintf("%v\n\n%v", verticesDesc, edgeDesc)
}

func (matrix *AdjacencyMatrix) BreadthFirstSearch(source Vertex) (visited []Vertex) {
	queue := Collections.NewQueueStack()
	enqueued := map[Vertex]bool{}

	queue.Enqueue(&source)
	enqueued[source] = true
	vertex := queue.Dequeue()
	for vertex != nil {
		obj := vertex.(*Vertex)
		visited = append(visited, *obj)
		for _, edge := range matrix.Edges(*obj) {
			_, ok := enqueued[edge.Destination]
			if !ok {
				queue.Enqueue(&edge.Destination)
				enqueued[edge.Destination] = true
			}
		}
		vertex = queue.Dequeue()
	}
	return
}
