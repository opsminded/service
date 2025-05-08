package service

import (
	"errors"
	"time"

	"github.com/opsminded/graphlib"
)

type Edge struct {
	Label       string
	Source      string
	Destination string
}

type Vertex struct {
	Label string
}

type SubGraph struct {
	Title      string
	Principal  Vertex
	Highlights []*Vertex
	Edges      []*Edge
	Vertices   []*Vertex
}

type Extractor interface {
	Frequency() time.Duration

	NextEdge() Edge
	HasNextEdge() bool

	NextVertex() Vertex
	HasNextVertex() bool
}

type Service struct {
	graph      graphlib.Graph
	extractors []Extractor
}

func New() *Service {
	return &Service{
		graph:      *graphlib.NewGraph(),
		extractors: []Extractor{},
	}
}

func (s *Service) GetVertex(label string) (Vertex, error) {
	v := s.graph.GetVertexByLabel(label)
	if v == nil {
		return Vertex{}, errors.New("vertex not found")
	}

	r := Vertex{
		Label: v.Label(),
	}

	return r, nil
}
