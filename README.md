# Trip Planner: A Graph-Based Travel Optimization Tool

This Go application helps you plan and optimize your travel route by calculating the most efficient path between different cities. It uses a weighted graph to determine the best route based on ticket prices, distance. 

## Features
- **Read CSV**: Import a CSV file containing travel data (cities, ticket prices, distances, and travel times).
- **Graph Representation**: Each city is represented as a node in a graph, with weighted edges indicating travel cost, distance.
- **Shortest Path Calculation**: Calculate the shortest route to a specified destination using Dijkstra’s algorithm.
- **Customization**: Adjust the weights for ticket price, distance, and time based on your preferences. (Not implemented)

## CSV File Format
The CSV file should follow this format, with columns for city ID, city name, destination city ID, ticket price, distance, and hours of travel:

```
id,city,to,ticket_average,distance,hours
0,Braga,0,0,0,0
1,Lisbon,0,40,50,5
2,Porto,0,25,30,1.5
2,Porto,1,30,26,1.5
```

### Column Descriptions:
- `id`: Unique ID for each city.
- `city`: Name of the city.
- `to`: Destination city ID.
- `ticket_average`: Average ticket price to the destination city.
- `distance`: Distance in kilometers to the destination city.
- `hours`: Travel time in hours to the destination city.

## Usage

### Command Line Flags:
- `-filepath`: Specify the path to your CSV file (default: `"example.csv"`).
- `-to`: Specify the ID of the desired destination city.
- `-help`: Display usage instructions.

### Example Command:
```bash
go run main.go -filepath="your_data.csv" -to=1
```

### Expected Output:
1. A graph visualization of the cities and their connections.
2. The shortest path from the start city (city ID `0`) to the specified destination city (based on ticket price, distance).
3. The travel path and the total weight (combined score of ticket price, distance, and hours).

## Code Breakdown

### 1. **CsvFile Struct**  
Defines the structure to hold data for each travel route: city ID, city name, destination city ID, ticket price, distance, and travel time.

### 2. **Graph Struct**  
Represents a graph of cities and travel routes:
- **CreateNodes()**: Adds nodes to the graph.
- **CreateEdges()**: Adds weighted edges between cities based on the calculated weight (ticket price, distance, and time).
- **GetShortestTo()**: Calculates the shortest path to a specified destination using Dijkstra’s algorithm.
- **String()**: Displays the graph as a string.
- **DisplayShortest()**: Displays the shortest path and its weight.

### 3. **Normalization**  
The `normalize()` function scales values (ticket prices, distances, and hours) between 0 and 1 to ensure fair comparison across different metrics.

### 4. **Weights Calculation**  
The `CalculateWeights()` function computes a weighted score for each route based on your preferences. The default weights are:
- 40% for ticket price.
- 60% for distance.

### 5. **Reading CSV**  
The `readCsv()` function parses the CSV file and returns a list of `CsvFile` objects.

### 6. **Main Execution**  
- Parse command-line flags.
- Read the CSV file.
- Build the graph.
- Calculate and display the shortest path to the specified destination.

## Example Output

Assuming a CSV with the following data:

```
id,city,to,ticket_average,distance,hours
0,Braga,0,0,0,0
1,Lisbon,0,40,50,5
2,Porto,0,25,30,1.5
2,Porto,1,30,26,1.5
```

You might run the program like this:

```bash
go run main.go -filepath="example.csv" -to=1
```

And get the following output:

```
====== Graph ======
Braga -> Porto
Porto -> Lisbon
Porto -> Porto
====== Graph ======
Shortest Path:
Braga -> Porto -> Lisbon (Weight: 0.34)
```

## Dependencies
- **gonum/graph**: For graph data structures and pathfinding (Dijkstra algorithm).

Install the required package:
```bash
go get gonum.org/v1/gonum/graph
```

## License
This project is open-source and released under the [MIT license](LICENSE).

## Author
Victor Milhomem

---
