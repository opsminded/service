package service

import (
	"context"
	"log"
	"time"

	"github.com/opsminded/graphlib"
)

type Clock interface {
	Now() time.Time
}

type DefaultClock struct{}

func (c DefaultClock) Now() time.Time {
	return time.Now()
}

type QueryResult struct {
	Title     string
	All       bool
	Principal graphlib.Vertex
	SubGraph  graphlib.Subgraph
}

type Summary struct {
	TotalVertex    int
	TotalEdges     int
	UnhealthVertex []graphlib.Vertex
}

type Extractor interface {
	Frequency() time.Duration

	NextEdge() graphlib.Edge
	HasNextEdge() bool

	NextVertex() graphlib.Vertex
	HasNextVertex() bool

	Reset()
}

type Service struct {
	clock         Clock
	extractors    []Extractor
	checkInterval time.Duration
	graph         graphlib.Graph
}

func New(ctx context.Context, checkInterval time.Duration, extractors []Extractor, clock Clock) *Service {
	if clock == nil {
		clock = DefaultClock{}
	}
	service := &Service{
		clock:         clock,
		extractors:    extractors,
		checkInterval: checkInterval,
		graph:         *graphlib.NewGraph(),
	}
	service.startExtractLoop(ctx)
	service.startHealthLoop(ctx)
	return service
}

func (s *Service) startExtractLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.extract()
			case <-ctx.Done():
				log.Println("extract loop done")
				return
			}
		}
	}()
}

func (s *Service) startHealthLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.CheckAndPropagate()
			case <-ctx.Done():
				log.Println("health loop done")
				return
			}
		}
	}()
}

func (s *Service) CheckAndPropagate() {

}

func (s *Service) extract() {
	for _, ex := range s.extractors {
		ex.Reset()

		for ex.HasNextVertex() {
			log.Println("next vertex")
			v := ex.NextVertex()
			s.graph.NewVertex(v.Label)
		}

		for ex.HasNextEdge() {
			log.Println("next edge")
			e := ex.NextEdge()
			s.graph.NewEdge(e.Label, e.Source.Label, e.Destination.Label)
		}
	}
}

func (s *Service) GetVertex(label string) (graphlib.Vertex, error) {
	v := s.graph.GetVertexByLabel(label)
	return v, nil
}

func (s *Service) Summary() Summary {
	sum := Summary{
		TotalVertex: s.graph.VertexLen(),
		TotalEdges:  s.graph.EdgeLen(),
	}
	return sum
}

func (s *Service) GetVertexDependencies(label string, all bool) QueryResult {
	dep := s.graph.GetVertexDependencies(label, all)
	sub := QueryResult{
		Title: "Dependências de " + label,
		All:   all,
		Principal: graphlib.Vertex{
			Label: label,
		},
		SubGraph: dep,
	}
	return sub
}

func (s *Service) GetVertexDependents(label string, all bool) QueryResult {
	dep := s.graph.GetVertexDependents(label, all)
	sub := QueryResult{
		Title: "Dependentes de " + label,
		All:   all,
		Principal: graphlib.Vertex{
			Label: label,
		},
		SubGraph: dep,
	}

	return sub
}

func (s *Service) Neighbors(label string) QueryResult {
	neighbors := s.graph.Neighbors(label)
	sub := QueryResult{
		Title:     "Vizinhos de " + label,
		Principal: graphlib.Vertex{Label: label},
		SubGraph:  neighbors,
	}
	return sub
}

func (s *Service) Path(label, destination string) QueryResult {
	path := s.graph.Path(label, destination)
	sub := QueryResult{
		Title: "Caminhos de " + label + " para " + destination,
		Principal: graphlib.Vertex{
			Label: label,
		},
		SubGraph: path,
	}
	return sub
}
