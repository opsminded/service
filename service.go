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
	Principal graphlib.Vertex
	SubGraph  graphlib.Subgraph
}

type Summary struct {
	TotalEdges        int
	TotalVertices     int
	UnhealthyVertices []graphlib.Vertex
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
		graph:         *graphlib.NewGraph(ctx),
	}
	service.startExtractLoop(ctx)
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

func (s *Service) extract() {
	for _, ex := range s.extractors {
		ex.Reset()
		for ex.HasNextVertex() {
			v := ex.NextVertex()
			s.graph.NewVertex(v.Key, v.Label, v.Healthy)
			log.Println("Vertex:", v.Key, v.Label, v.Healthy)
		}
		for ex.HasNextEdge() {
			e := ex.NextEdge()
			s.graph.NewEdge(e.Source, e.Target)
			log.Println("Edge:", e.Source, e.Target)
		}
	}
}

func (s *Service) GetVertex(key string) (graphlib.Vertex, error) {
	return s.graph.GetVertex(key)
}

func (s *Service) SetVertexHealth(key string, health bool) error {
	return s.graph.SetVertexHealth(key, health)
}

func (s *Service) ClearGraphHealthyStatus() {
	s.graph.ClearGraphHealthyStatus()
}

func (s *Service) Summary() Summary {
	// TODO: adicionar os contadores de dependências
	sum := Summary{
		TotalEdges:        0,
		TotalVertices:     0,
		UnhealthyVertices: nil,
	}
	return sum
}

func (s *Service) VertexDependencies(key string, all bool) (QueryResult, error) {
	p, err := s.graph.GetVertex(key)
	if err != nil {
		return QueryResult{}, err
	}

	dep, err := s.graph.VertexDependencies(key, all)
	if err != nil {
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Dependências de " + p.Label,
		Principal: p,
		SubGraph:  dep,
	}
	return sub, nil
}

func (s *Service) GetVertexDependents(key string, all bool) (QueryResult, error) {
	p, err := s.graph.GetVertex(key)
	if err != nil {
		return QueryResult{}, err
	}

	dep, err := s.graph.VertexDependents(key, all)
	if err != nil {
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Dependentes de " + p.Label,
		Principal: p,
		SubGraph:  dep,
	}
	return sub, nil
}

func (s *Service) Neighbors(key string) (QueryResult, error) {
	p, err := s.graph.GetVertex(key)
	if err != nil {
		return QueryResult{}, err
	}

	// neighbors, err := s.graph.Neighbors(label)
	// if err != nil {
	// 	return QueryResult{}, err
	// }

	sub := QueryResult{
		Title:     "Vizinhos de " + p.Label,
		Principal: p,
		SubGraph:  graphlib.Subgraph{},
	}
	return sub, nil
}

func (s *Service) Path(kSrc, ktgt string) (QueryResult, error) {
	src, err := s.graph.GetVertex(kSrc)
	if err != nil {
		return QueryResult{}, err
	}

	tgt, err := s.graph.GetVertex(ktgt)
	if err != nil {
		return QueryResult{}, err
	}

	path, err := s.graph.Path(kSrc, ktgt)
	if err != nil {
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Caminhos de " + src.Label + " para " + tgt.Label,
		Principal: src,
		SubGraph:  path,
	}
	return sub, nil
}

func (s *Service) GetVertexLineages(key string) (QueryResult, error) {
	p, err := s.graph.GetVertex(key)
	if err != nil {
		return QueryResult{}, err
	}

	// lineage, err := s.graph.Lineage(label)
	// if err != nil {
	// 	return QueryResult{}, err
	// }

	sub := QueryResult{
		Title:     "Linhagens de " + p.Label,
		Principal: p,
		SubGraph:  graphlib.Subgraph{},
	}

	return sub, nil
}
