package Graph

import (
	"GiveMeAnOfferGo/Collections/Graph"
	"GiveMeAnOfferGo/test/Utils"
	"fmt"
	"testing"
)

func TestAdjacencyList(t *testing.T) {
	graph := Graph.NewAdjacencyList()
	singapore := graph.CreateVertex(Utils.GetString("Singapore"))
	tokyo := graph.CreateVertex(Utils.GetString("Tokyo"))
	hongKong := graph.CreateVertex(Utils.GetString("Hong Kong"))
	detroit := graph.CreateVertex(Utils.GetString("Detroit"))
	sanFrancisco := graph.CreateVertex(Utils.GetString("San Francisco"))
	washingtonDC := graph.CreateVertex(Utils.GetString("Washington DC"))
	austinTexas := graph.CreateVertex(Utils.GetString("Austin Texas"))
	seattle := graph.CreateVertex(Utils.GetString("Seattle"))

	graph.Add(Graph.EdgeTypeUndirected, singapore, hongKong, 300)
	graph.Add(Graph.EdgeTypeUndirected, singapore, tokyo, 500)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, tokyo, 250)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, detroit, 450)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, washingtonDC, 300)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, sanFrancisco, 600)
	graph.Add(Graph.EdgeTypeUndirected, detroit, austinTexas, 50)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, washingtonDC, 292)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, washingtonDC, 337)
	graph.Add(Graph.EdgeTypeUndirected, washingtonDC, seattle, 277)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, seattle, 218)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, sanFrancisco, 297)

	t.Log(graph)

	fmt.Println("San Francisco Outgoing Flights:")
	fmt.Println("-------------------------------")
	for _, edge := range graph.Edges(sanFrancisco) {
		fmt.Printf("from: %v to %v\n", edge.Source, edge.Destination)
	}
}

func TestAdjacencyMatrix(t *testing.T) {
	graph := Graph.NewAdjacencyMatrix()
	singapore := graph.CreateVertex(Utils.GetString("Singapore"))
	tokyo := graph.CreateVertex(Utils.GetString("Tokyo"))
	hongKong := graph.CreateVertex(Utils.GetString("Hong Kong"))
	detroit := graph.CreateVertex(Utils.GetString("Detroit"))
	sanFrancisco := graph.CreateVertex(Utils.GetString("San Francisco"))
	washingtonDC := graph.CreateVertex(Utils.GetString("Washington DC"))
	austinTexas := graph.CreateVertex(Utils.GetString("Austin Texas"))
	seattle := graph.CreateVertex(Utils.GetString("Seattle"))

	graph.Add(Graph.EdgeTypeUndirected, singapore, hongKong, 300)
	graph.Add(Graph.EdgeTypeUndirected, singapore, tokyo, 500)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, tokyo, 250)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, detroit, 450)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, washingtonDC, 300)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, sanFrancisco, 600)
	graph.Add(Graph.EdgeTypeUndirected, detroit, austinTexas, 50)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, washingtonDC, 292)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, washingtonDC, 337)
	graph.Add(Graph.EdgeTypeUndirected, washingtonDC, seattle, 277)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, seattle, 218)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, sanFrancisco, 297)

	t.Log(graph)

	fmt.Println("San Francisco Outgoing Flights:")
	fmt.Println("-------------------------------")
	for _, edge := range graph.Edges(sanFrancisco) {
		fmt.Printf("from: %v to %v\n", edge.Source, edge.Destination)
	}
}

func TestBFS(t *testing.T) {
	graph := Graph.NewAdjacencyList()
	singapore := graph.CreateVertex(Utils.GetString("Singapore"))
	tokyo := graph.CreateVertex(Utils.GetString("Tokyo"))
	hongKong := graph.CreateVertex(Utils.GetString("Hong Kong"))
	detroit := graph.CreateVertex(Utils.GetString("Detroit"))
	sanFrancisco := graph.CreateVertex(Utils.GetString("San Francisco"))
	washingtonDC := graph.CreateVertex(Utils.GetString("Washington DC"))
	austinTexas := graph.CreateVertex(Utils.GetString("Austin Texas"))
	seattle := graph.CreateVertex(Utils.GetString("Seattle"))

	graph.Add(Graph.EdgeTypeUndirected, singapore, hongKong, 300)
	graph.Add(Graph.EdgeTypeUndirected, singapore, tokyo, 500)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, tokyo, 250)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, detroit, 450)
	graph.Add(Graph.EdgeTypeUndirected, tokyo, washingtonDC, 300)
	graph.Add(Graph.EdgeTypeUndirected, hongKong, sanFrancisco, 600)
	graph.Add(Graph.EdgeTypeUndirected, detroit, austinTexas, 50)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, washingtonDC, 292)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, washingtonDC, 337)
	graph.Add(Graph.EdgeTypeUndirected, washingtonDC, seattle, 277)
	graph.Add(Graph.EdgeTypeUndirected, sanFrancisco, seattle, 218)
	graph.Add(Graph.EdgeTypeUndirected, austinTexas, sanFrancisco, 297)

	for _, vertex := range graph.BreadthFirstSearch(austinTexas) {
		fmt.Println(vertex)
	}
}
