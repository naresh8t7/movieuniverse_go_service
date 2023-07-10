# MovieUniverse Application

### Stack

* Go
* grpc-ecosystem/grpc-gateway
* swagger for REST API documentation
* Testify
* https://github.com/vektra/mockery
* gRPC, Protobufs
* [Golang Neo4j SeaBolt Driver](https://github.com/neo4j-drivers/seabolt/releases)
* Neo4j-Server
* Frontend: React, Redux, bootstrap, [d3.js](http://d3js.org/) The frontend instructions are provided in the [repo](https://github.com/naresh8t7/movie_uiniverse_react_app) .

### Architecture 

![Movie Universe](https://github.com/naresh8t7/movieuniverse_go_service/blob/main/Movie%20Universe%20Architecture.png)

### Endpoints:

Get Movie

### Setup

This uses the Go standard library HTTP server, along with the Golang Neo4j Bolt Driver library.

### Run locally:

Start your local Neo4j Server [Download & Install](http://neo4j.com/download), open the [Neo4j Browser](http://localhost:7474).
Then install the create movie db based on movies.schema by running following command, and hit the triangular "Run" button.
```bash
CALL graphql.idl('type Movie  {
  title: String!
  released: Int
  actors: [Person] @relation(name:"ACTED_IN",direction:IN)
  directors: [Person] @cypher(statement:"MATCH (this)<-[:DIRECTED]-(d) RETURN d")
  tags: [String]
}
type Person {
  name: String!
  born: Int
  movies: [Movie] @relation(name:"ACTED_IN")
}
schema {
   query: QueryType
   mutation: MutationType
}
type QueryType {
  coActors(name:ID!): [Person] @cypher(statement:"MATCH (p:Person {name:$name})-[:ACTED_IN]->()<-[:ACTED_IN]-(co) RETURN distinct co")
}
type MutationType {
  rateMovie(user:ID!, movie:ID!, rating:Int!): Int
  @cypher(statement:"MATCH (p:Person {name:$user}),(m:Movie {title:$movie}) MERGE (p)-[r:RATED]->(m) SET r.rating=$rating RETURN r.rating")
}')
```
To build the application.
```
go build ./...
```

To test the units test.
```
go test ./...
```
Start this application with:

```bash
set the pwd to movieuniverse/movie/cmd/serve and run below command.
server.exe -grpc-port=9090 -db-host=bolt://localhost:7687 -db-user=neo4j -db-password=test --http-port=8080
```

Use below curl commands to verify the get movie and list all movies.
```bash
# JSON object for single movie with cast
curl http://localhost:8080/v1/movie/The%20Matrix

# list of JSON objects for all movie search results
curl http://localhost:8080/v1/movie/all
```

Also included client-grpc and client-rest to verify the services developed in the application.
```bash
set the pwd to movieuniverse\movie\cmd\client-grpc and run below command.
client-grpc.exe -server=localhost:9090

set the pwd to movieuniverse\movie\cmd\client-rest and run below command.
client-rest.exe -server=localhost:8080
```
