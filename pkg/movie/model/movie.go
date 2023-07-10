package model

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	v1 "movieuniverse/pkg/movie/api/proto/v1"
	"sync"
)

func CreateMovie(movie *v1.Movie, session neo4j.Session, ce chan error) {
	c := make(chan error)
	log.Printf("movie: %v", movie)
	// creating a movie
	result, err := session.Run(`CREATE (n:Movie {title:$title, tags:$tags, 
		released: $released}) 
		RETURN n.title`, map[string]interface{}{
		"title":                  movie.Title,
		"tags":                movie.Tags,
		"released":              movie.Year,
	})
	if err != nil {
		ce <- err
		return
	}
	result.Next()
	log.Println(result.Record().GetByIndex(0).(string))
	if err = result.Err(); err != nil {
		ce <- err
		return
	}

	// create Actor and Director nodes whenever an Movie is created
	var mutex = &sync.Mutex{}
	for _, actor := range movie.Actors {
		go createPersons(session, actor, movie.Title, "ACTED_IN", c, mutex)
		if err1 :=  <-c; err1 != nil {
			log.Printf("create actors failed %v", err1)
			return
		}

	}
	for _, director := range movie.Directors {
		go createPersons(session, director, movie.Title, "DIRECTED_BY", c, mutex)
		if err1 :=  <-c; err1 != nil {
			log.Printf("create actors failed %v", err1)
			return
		}
	}

	log.Println("Created Movie node")
	ce <- nil
	return
}

func AddTags(session  neo4j.Session, req *v1.AddTagsRequest, resp *v1.AddTagsResponse) error {
	result, err := session.Run(`MATCH (n:Movie{title: $title }) SET n.tags = $tags 
		RETURN n.title as title, n.year as year, n.tags as tags`, map[string]interface{}{
		"title":                  req.Title,
		"tags": req.Tags,
	})
	if err != nil {
		return err
	}
	movie := &v1.Movie{}
	for result.Next() {

		title, ok := result.Record().GetByIndex(0).(string)  //n.id
		if ok {
			movie.Title = title
		}
		year, ok := result.Record().GetByIndex(1).(int64)  //n.id
		if ok {
			movie.Year = year
		}
		tags, ok := result.Record().GetByIndex(2).([]string)  //n.id
		if ok {
			movie.Tags = tags
		}
	}
	resp.Movie = movie
	return nil
}

func ListMovies(session  neo4j.Session, req *v1.ListAllMoviesRequest, resp *v1.ListAllMoviesResponse) error {
	result, err := session.Run(`MATCH (n:Movie) 
		RETURN n.title as title, n.year as year, n.tags as tags`, map[string]interface{}{})
	if err != nil {
		return err
	}
	for result.Next() {
		movie := &v1.Movie{}
		title, ok := result.Record().GetByIndex(0).(string)  //n.id
		if ok {
			movie.Title = title
		}
		year, ok := result.Record().GetByIndex(1).(int64)  //n.id
		if ok {
			movie.Year = year
		}
		tags, ok := result.Record().GetByIndex(2).([]string)  //n.id
		if ok {
			movie.Tags = tags
		}
		resp.Movies = append(resp.Movies, movie)
	}
	return nil
}

func GetMovie(session  neo4j.Session, req *v1.GetMovieRequest, resp *v1.GetMovieResponse) error {
	result, err := session.Run(`MATCH (n:Movie) where n.title = $title
		RETURN n.title as title, n.year as year, n.tags as tags`, map[string]interface{}{
		"title":                  req.Title,
	})
	if err != nil {
		return err
	}
	movie := &v1.Movie{}
	for result.Next() {

		title, ok := result.Record().GetByIndex(0).(string)  //n.id
		if ok {
			movie.Title = title
		}
		year, ok := result.Record().GetByIndex(1).(int64)  //n.id
		if ok {
			movie.Year = year
		}
		tags, ok := result.Record().GetByIndex(2).([]string)  //n.id
		if ok {
			movie.Tags = tags
		}
	}
	resp.Movie = movie
	return nil
}

func createPersons(session neo4j.Session, person *v1.Person, title string, label string, c chan error, mutex *sync.Mutex) {
	// beginning of the critical section
	mutex.Lock()
	result, err := session.Run(`MATCH(a:Movie) WHERE a.title=$title
// low level database functions
	CREATE (n:Actor {name:$name})->[:`+label+`]-(a) `, map[string]interface{}{
		"title":          title,
		"name":           person.Name,
	})
	if err != nil {
		c <- err
		return
	}
	// critical section ends
	mutex.Unlock()
	if err = result.Err(); err != nil {
		c <- err
		return
	}
	log.Printf("Created %s node", label)
	c <- nil
	return
}
