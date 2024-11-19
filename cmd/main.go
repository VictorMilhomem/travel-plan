package main

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

/*
id -> city ID int64
city -> city name string
to -> city ID int64
ticket_average -> ticket price from City -> To columns float32
distance -> distance KM from City -> To columns -> int64
hours -> hours to go from City -> To

id, city, to, ticket_average, distance, hours
0, Braga, 0, 0, 0, 0
1, Lisbon, 0,40, 50, 5
2, Porto, 0, 25, 30, 1.5
2, Porto, 1, 30, 26, 1.5


*/

type CsvFile struct {
	ID             int64
	city           string
	to             int64
	ticket_average float32
	distance       int64
	hours          float32
}

func NewCsvFile(
	id int64,
	city string,
	to int64,
	ticket_average float32,
	distance int64,
	hours float32,
) *CsvFile {
	return &CsvFile{
		id,
		city,
		to,
		ticket_average,
		distance,
		hours,
	}
}

func (self *CsvFile) String() string {
	return fmt.Sprintf("id %d, city %s, to %d, ticket_average %f, distance %d, hours %f", self.ID, self.city, self.to, self.ticket_average, self.distance, self.hours)
}

type Graph struct {
	data      []*CsvFile
	nodes_ids map[int64]string
	g         *simple.WeightedUndirectedGraph
}

func NewGraph(data []*CsvFile) *Graph {
	node_ids := make(map[int64]string)
	// get all different nodes
	for i := range data {
		id := data[i].ID
		city := data[i].city
		node_ids[id] = city
	}

	return &Graph{
		data:      data,
		nodes_ids: node_ids,
		g:         simple.NewWeightedUndirectedGraph(0, 0),
	}
}

func (self *Graph) CreateNodes() {
	// create all the nodes
	for key := range self.nodes_ids {
		node := simple.Node(key)
		self.g.AddNode(node)
	}
}

func (self *Graph) CreateEdges() {
	// create the edges
	for i := range self.data {
		from := self.g.Node(self.data[i].ID)
		to := self.g.Node(self.data[i].to)
		if from != nil && to != nil && from != to {
			// create a weight
			self.g.SetWeightedEdge(self.g.NewWeightedEdge(from, to, 3.0))
		}
	}
}

func (self *Graph) GetShortestTo(from, to int64) ([]graph.Node, float64) {
	shortest := path.DijkstraFrom(simple.Node(from), self.g)
	_path, weight := shortest.To(simple.Node(to).ID())
	return _path, weight
}

func (self *Graph) String() string {
	var builder strings.Builder

	for edges := self.g.Edges(); edges.Next(); {
		from := edges.Edge().From().ID()
		to := edges.Edge().To().ID()
		builder.WriteString(self.nodes_ids[from] + " -> " + self.nodes_ids[to])
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (self *Graph) DisplayShortest(_path []graph.Node, weight float64) {
	var builder strings.Builder

	for i, node := range _path {
		nodeID := node.ID()
		if i == len(_path)-1 {
			builder.WriteString(self.nodes_ids[nodeID])
		} else {
			builder.WriteString(self.nodes_ids[nodeID] + " -> ")
		}
	}

	// Append weight at the end if relevant
	builder.WriteString(fmt.Sprintf(" (Weight: %.2f)", weight))
	fmt.Println(builder.String())
}

func main() {
	var data []*CsvFile

	data = append(data, NewCsvFile(0, "Braga", 0, 0, 0, 0))
	data = append(data, NewCsvFile(1, "Lisbon", 0, 40, 50, 5))
	data = append(data, NewCsvFile(2, "Porto", 0, 25, 30, 1.5))
	data = append(data, NewCsvFile(2, "Porto", 1, 30, 26, 1.5))

	g := NewGraph(data)
	g.CreateNodes()
	g.CreateEdges()
	fmt.Println("====== Graph ======")
	fmt.Println(g.String())

	_path, weight := g.GetShortestTo(0, 2)
	g.DisplayShortest(_path, weight)
}