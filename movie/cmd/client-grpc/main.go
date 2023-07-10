package main

import (
	v1 "movieuniverse/pkg/movie/api/proto/v1"
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewMovieServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	var actors []*v1.Person
	actor := &v1.Person{
		Name:                 "Tom Cruise",
	}
	actors = append(actors, actor)
	// Call CreateMovie
	req1 := v1.CreateMovieRequest{
		Movie: &v1.Movie{
			Title:       "MI: 2",
			Year: 1989,
			Actors: actors,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("CreateMovie failed: %v", err)
	}
	log.Printf("CreateMovie result: <%+v>\n\n", res1)


	// GetMovie
	req2 := v1.GetMovieRequest{
		Title:  "MI: 2",
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("GetMovie failed: %v", err)
	}
	log.Printf("GetMovie result: <%+v>\n\n", res2)

	// Call GetMovieAll
	req3 := v1.ListAllMoviesRequest{
	}
	res3, err := c.ListAll(ctx, &req3)
	if err != nil {
		log.Fatalf("GetMovieAll failed: %v", err)
	}
	log.Printf("GetMovieAll result: <%+v>\n\n", res3)

	// Add tags
	req4 := v1.AddTagsRequest{
		Title:                "MI: 2",
		Tags:                 []string{"Action"},
	}
	res4, err := c.AddTags(ctx, &req4)
	if err != nil {
		log.Fatalf("AddTags failed: %v", err)
	}
	log.Printf("AddTags result: <%+v>\n\n", res4)
}
