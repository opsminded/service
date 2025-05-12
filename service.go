package service

import (
	"log"
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
	All        bool
	Highlights []Vertex
	Edges      []Edge
	Vertices   []Vertex
}

type Summary struct {
	TotalVertex    int
	TotalEdges     int
	UnhealthVertex []Vertex
}

type Extractor interface {
	Frequency() time.Duration

	NextEdge() Edge
	HasNextEdge() bool

	NextVertex() Vertex
	HasNextVertex() bool

	Reset()
}

type Service struct {
	graph      graphlib.Graph
	extractors []Extractor
}

func New(extractors []Extractor) *Service {
	return &Service{
		graph:      *graphlib.NewGraph(),
		extractors: extractors,
	}
}

func (s *Service) Extract() {
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
			s.graph.NewEdge(e.Label, e.Source, e.Destination)
		}
	}
}

func (s *Service) GetVertex(label string) (Vertex, error) {
	v := s.graph.GetVertexByLabel(label)
	r := Vertex{
		Label: v.Label,
	}
	return r, nil
}

func (s *Service) Summary() Summary {
	sum := Summary{
		TotalVertex: s.graph.VertexLen(),
		TotalEdges:  s.graph.EdgeLen(),
	}
	return sum
}

func (s *Service) GetVertexDependencies(label string, all bool) SubGraph {
	res := s.graph.GetVertexDependencies(label, all)
	sub := SubGraph{
		Title: "Vizinhos de " + label,
		All:   all,
		Principal: Vertex{
			Label: label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}
	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label,
		})
	}

	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		})
	}
	return sub
}

func (s *Service) GetVertexDependents(label string, all bool) SubGraph {
	res := s.graph.GetVertexDependents(label, all)

	sub := SubGraph{
		Title: "Dependencias de " + label,
		All:   all,
		Principal: Vertex{
			Label: label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label,
		})
	}

	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		})
	}
	return sub
}

func (s *Service) Neighbors(label string) SubGraph {

	principal := s.graph.GetVertexByLabel(label)

	sub := SubGraph{
		Title: "Vizinhos de " + label,
		Principal: Vertex{
			Label: principal.Label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	rs := s.graph.Neighbors(label)

	for _, v := range rs.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label,
		})
	}

	for _, e := range rs.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		})
	}

	return sub
}

func (s *Service) Path(label, destination string) SubGraph {
	res := s.graph.Path(label, destination)

	sub := SubGraph{
		Title: "Caminho de " + label + " para " + destination,
		Principal: Vertex{
			Label: label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label,
		})
	}
	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		})
	}
	return sub
}
