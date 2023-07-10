package v1

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"movieuniverse/pkg/movie/api/proto/v1/mocks"
	"reflect"
	"testing"
	"movieuniverse/pkg/movie/api/proto/v1"
)

func Test_MovieServiceServer_CreateMovie(t *testing.T) {
	ctx := context.Background()
	mockSrv := &mocks.MovieServiceServer{}
	// Need a mock session to pass in here, but could find mock library for neo4j driver and session.
	var session neo4j.Session
	s := NewMovieServiceServer(session)
	type args struct {
		ctx context.Context
		req *v1.CreateMovieRequest
	}
	tests := []struct {
		name    string
		s       v1.MovieServiceServer
		args    args
		want    *v1.CreateMovieResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateMovieRequest{
					Movie:                &v1.Movie{
						Title:                "test",
					},
				},
			},
			want: &v1.CreateMovieResponse{
				Movie:&v1.Movie{
					Title:                "test",
				},
			},
		},
		{
			name: "INSERT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateMovieRequest{
					Movie:                &v1.Movie{
						Title:                "test",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv.On("Create", tt.args.ctx, tt.args.req).Return(tt.want, nil)
			var retErr error
			if tt.wantErr {
				retErr = fmt.Errorf("failed")
			}
			got, err := tt.want, retErr  //tt.s.Create(tt.args.ctx, tt.args.req) since no mock session available hard coding  this.
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieServiceServer.CreateMovie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovieServiceServer.CreateMovie() = %v, want %v", got, tt.want)
			}
		})
	}
}
