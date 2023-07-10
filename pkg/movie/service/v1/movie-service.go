package v1

import (
	v1 "movieuniverse/pkg/movie/api/proto/v1"
	model "movieuniverse/pkg/movie/model"
	"context"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// movieServiceServer is implementation of v1.MovieServiceServer proto interface
type movieServiceServer struct {
	session neo4j.Session
}

func (m *movieServiceServer) AddTags(ctx context.Context, req *v1.AddTagsRequest) (*v1.AddTagsResponse, error) {
	if req.Title == "" || len(req.Tags) == 0 {
		return nil, status.Error(codes.Internal, "Required fields missing")
	}
	resp :=  &v1.AddTagsResponse{}
	err := model.AddTags(m.session, req, resp)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to add tags  movie %v ",err)
	}
	return resp, nil
}

func (m *movieServiceServer) ListAll(ctx context.Context, req *v1.ListAllMoviesRequest) (*v1.ListAllMoviesResponse, error) {
	if req.Filter != "" {
		// TODO(Naresh): Implement movie filter.
		log.Println("Movie filter is not supported yet.")
	}
	resp :=  &v1.ListAllMoviesResponse{}
	err := model.ListMovies(m.session, req, resp)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to list  movie %v ",err)
	}
	return resp, nil
}

func (m *movieServiceServer) Create(ctx context.Context, req *v1.CreateMovieRequest) (*v1.CreateMovieResponse, error) {
	// insert movie entity data
	ce := make(chan error)
	// goroutine for invoking the model layer event create function
	go model.CreateMovie(req.GetMovie(), m.session,  ce)
	if err := <-ce; err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to create movie-> "+err.Error())
	}
	resp :=  &v1.CreateMovieResponse{
		Movie:req.Movie,
	}
	return resp, nil
}

func (m *movieServiceServer) Read(ctx context.Context, req *v1.GetMovieRequest) (*v1.GetMovieResponse, error) {
	resp :=  &v1.GetMovieResponse{}
	err := model.GetMovie(m.session, req, resp)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get  movie %v ",err)
	}
	return resp, nil
}

// NewMovieServiceServer creates Movie service
func NewMovieServiceServer(session neo4j.Session) v1.MovieServiceServer {
	return &movieServiceServer{session:session}
}
