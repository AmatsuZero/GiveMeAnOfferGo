package Graph

import (
	"GiveMeAnOfferGo/Objects"
	"fmt"
)

type EdgeType int

const (
	EdgeTypeDirected EdgeType = iota
	EdgeTypeUndirected
)

type Vertex struct {
	Index int
	Data  Objects.EquatableObject
}

func (vertex *Vertex) IsNil() bool {
	return vertex.Data == nil
}

func (vertex *Vertex) IsEqualTo(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if obj == vertex {
		return true
	}
	object, ok := obj.(Vertex)
	if !ok {
		tmp, ok := obj.(*Vertex)
		if ok {
			object = *tmp
		} else {
			return false
		}
	}
	return vertex.Index == object.Index && vertex.Data.IsEqualTo(object.Data)
}

func (vertex *Vertex) String() string {
	return fmt.Sprintf("%v: %v", vertex.Index, vertex.Data)
}

type Edge struct {
	Source      Vertex
	Destination Vertex
	Weight      float64
}

type Graph interface {
	CreateVertex(data Objects.EquatableObject) Vertex
	AddDirectedEdge(source Vertex, destination Vertex, weight float64)
	AddUndirectedEdge(source Vertex, destination Vertex, weight float64)
	Add(edgeType EdgeType, source Vertex, destination Vertex, weight float64)
	Edges(source Vertex) []Edge
	Weight(source Vertex, destination Vertex) float64
}
