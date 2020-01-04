package Graph

import (
	"GiveMeAnOfferGo/Collections"
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type AdjacencyList struct {
	adjacencies map[Vertex][]Edge
}

func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{adjacencies: map[Vertex][]Edge{}}
}

func (adjacencyList *AdjacencyList) CreateVertex(data Objects.EquatableObject) Vertex {
	vertex := Vertex{
		Index: len(adjacencyList.adjacencies),
		Data:  data,
	}
	adjacencyList.adjacencies[vertex] = make([]Edge, 0)
	return vertex
}

func (adjacencyList *AdjacencyList) AddDirectedEdge(source Vertex, destination Vertex, weight float64) {
	edge := Edge{
		Source:      source,
		Destination: destination,
		Weight:      weight,
	}
	adjacencyList.adjacencies[source] = append(adjacencyList.adjacencies[source], edge)
}

func (adjacencyList *AdjacencyList) AddUndirectedEdge(source Vertex, destination Vertex, weight float64) {
	adjacencyList.AddDirectedEdge(source, destination, weight)
	adjacencyList.AddDirectedEdge(destination, source, weight)
}

func (adjacencyList *AdjacencyList) Add(edgeType EdgeType, source Vertex, destination Vertex, weight float64) {
	switch edgeType {
	case EdgeTypeDirected:
		adjacencyList.AddDirectedEdge(source, destination, weight)
	case EdgeTypeUndirected:
		adjacencyList.AddUndirectedEdge(source, destination, weight)
	}
}

func (adjacencyList *AdjacencyList) Edges(source Vertex) []Edge {
	return adjacencyList.adjacencies[source]
}

func (adjacencyList *AdjacencyList) Weight(source Vertex, destination Vertex) (weight float64) {
	weight = Objects.InvalidResult
	for _, dest := range adjacencyList.Edges(source) {
		if dest.Destination.IsEqualTo(destination) {
			weight = dest.Weight
		}
	}
	return
}

func (adjacencyList *AdjacencyList) String() (result string) {
	result = ""
	for vertex, edges := range adjacencyList.adjacencies {
		edgeString := ""
		for index, edge := range edges {
			if index != len(edges)-1 {
				edgeString += fmt.Sprintf("%v, ", edge.Destination)
			} else {
				edgeString += fmt.Sprint(edge.Destination)
			}
		}
		result += fmt.Sprintf("%v ---> [ %v ]\n", vertex, edgeString)
	}
	return
}

func (adjacencyList *AdjacencyList) BreadthFirstSearch(source Vertex) (visited []Vertex) {
	queue := Collections.NewQueueStack()
	enqueued := map[Vertex]bool{}

	queue.Enqueue(&source)
	enqueued[source] = true
	vertex := queue.Dequeue()
	for vertex != nil {
		obj := vertex.(*Vertex)
		visited = append(visited, *obj)
		for _, edge := range adjacencyList.Edges(*obj) {
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
