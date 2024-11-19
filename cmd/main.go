package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

func normalize(vector []float32) []float32 {
	min := slices.Min(vector)
	max := slices.Max(vector)

	// Normalize the vector
	var ret []float32
	for _, value := range vector {
		normalizedValue := (value - min) / (max - min)
		ret = append(ret, normalizedValue)
	}

	return ret
}

func (self *Graph) CalculateWeights() []float32 {
	preferences := map[string]float32{
		"weight_ticket":   0.4,
		"weight_distance": 0.6,
	}
	var scores []float32
	// normalize the data
	var dist []float32
	var ticket []float32
	for i := range self.data {
		dist = append(dist, float32(self.data[i].distance))
		ticket = append(ticket, self.data[i].ticket_average)
	}

	norm_dist := normalize(dist)
	norm_ticket := normalize(ticket)

	// define the rule for the weight
	for i := range dist {
		score := (preferences["weight_ticket"] * norm_ticket[i]) + (preferences["weight_distance"] * norm_dist[i])
		scores = append(scores, score)
	}

	return scores
}

func (self *Graph) CreateEdges() {
	// create the edges
	weights := self.CalculateWeights()
	for i := range self.data {
		from := self.g.Node(self.data[i].ID)
		to := self.g.Node(self.data[i].to)
		if from != nil && to != nil && from != to {
			self.g.SetWeightedEdge(self.g.NewWeightedEdge(from, to, float64(weights[i])))
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

	builder.WriteString(fmt.Sprintf(" (Weight: %.2f)", weight))
	fmt.Println(builder.String())
}

func readCsv(path string) []*CsvFile {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open file %s\n", path)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading file:", err)
	}

	var data []*CsvFile
	for i, row := range records {
		if i == 0 {
			continue
		}
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			log.Printf("Error parsing ID at line %d: %v", i+1, err)
			continue
		}
		city := row[1]
		to, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			log.Printf("Error parsing To at line %d: %v", i+1, err)
			continue
		}
		ticketAverage, err := strconv.ParseFloat(row[3], 32)
		if err != nil {
			log.Printf("Error parsing Ticket Average at line %d: %v", i+1, err)
			continue
		}
		distance, err := strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			log.Printf("Error parsing Distance at line %d: %v", i+1, err)
			continue
		}
		hours, err := strconv.ParseFloat(row[5], 32)
		if err != nil {
			log.Printf("Error parsing Hours at line %d: %v", i+1, err)
			continue
		}

		data = append(data, NewCsvFile(id, city, to, float32(ticketAverage), int64(distance), float32(hours)))
	}
	return data
}

func main() {
	var data []*CsvFile

	path := flag.String("filepath", "example.csv", "Specify the filepath")
	to := flag.Int64("to", 1, "Specify the ID of the desired place to go")
	help := flag.Bool("help", false, "Display usage instructions")
	flag.Parse()

	if *help || len(os.Args) == 1 {
		fmt.Println("Usage:\tplan -filepath=<your_filepath>.csv -to=<destination_id>\n")
		fmt.Println("\t=============== Your CSV file should be in the following format ==================")
		fmt.Println("\tid,city,to,ticket_average,distance,hours")
		fmt.Println("\t0, Braga, 0, 0, 0, 0")
		fmt.Println("\t1, Lisbon, 0,40, 50, 5")
		fmt.Println("\t2, Porto, 0, 25, 30, 1.5")
		fmt.Println("\t2, Porto, 1, 30, 26, 1.5")
		return
	}

	data = readCsv(*path)
	g := NewGraph(data)
	g.CreateNodes()
	g.CreateEdges()

	fmt.Println("====== Graph ======")
	fmt.Println(g.String())
	fmt.Println("====== Graph ======")

	_path, weight := g.GetShortestTo(0, *to)
	g.DisplayShortest(_path, weight)
}
