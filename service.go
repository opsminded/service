package service

import (
	"errors"
	"log"
	"time"

	"github.com/opsminded/graphlib"
)

type Edge struct {
	Label       string
	Class       string
	Source      string
	Destination string
}

type Vertex struct {
	Label string
	Class string
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
			if _, err := s.graph.NewVertex(v.Label, v.Class); err != nil {
				panic(err)
			}
		}

		for ex.HasNextEdge() {
			log.Println("next edge")
			e := ex.NextEdge()
			v1 := s.graph.GetVertexByLabel(e.Source)
			v2 := s.graph.GetVertexByLabel(e.Destination)
			if _, err := s.graph.NewEdge(e.Label, e.Class, v1, v2); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Service) GetVertex(label string) (Vertex, error) {
	v := s.graph.GetVertexByLabel(label)
	if v == nil {
		return Vertex{}, errors.New("vertex not found")
	}

	r := Vertex{
		Label: v.Label(),
		Class: s.graph.GetVertexClass(v),
	}

	return r, nil
}

func (s *Service) Summary() Summary {
	sum := Summary{
		TotalVertex: s.graph.VertexLen(),
		TotalEdges:  s.graph.EdgeLen(),
	}

	list := s.graph.UnhealthVertices()
	for _, v := range list {
		sum.UnhealthVertex = append(sum.UnhealthVertex, Vertex{
			Label: v.Label(),
			Class: s.graph.GetVertexClass(v),
		})
	}

	return sum
}

func (s *Service) GetVertexDependencies(label string, all bool) SubGraph {

	res := s.graph.GetVertexDependencies(label, all)

	sub := SubGraph{
		Title: "Vizinhos de " + label,
		All:   res.All,
		Principal: Vertex{
			Label: res.Principal.Label(),
			Class: s.graph.GetVertexClass(res.Principal),
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label(),
			Class: s.graph.GetVertexClass(v),
		})
	}

	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label(),
			Class:       s.graph.GetEdgeClass(e),
			Source:      s.graph.EdgeSourceLabel(e),
			Destination: s.graph.EdgeDestinationLabel(e),
		})
	}
	return sub
}

func (s *Service) GetVertexDependants(label string, all bool) SubGraph {
	res := s.graph.GetVertexDependants(label, all)

	sub := SubGraph{
		Title: "Dependencias de " + label,
		All:   res.All,
		Principal: Vertex{
			Label: res.Principal.Label(),
			Class: s.graph.GetVertexClass(res.Principal),
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label(),
			Class: s.graph.GetVertexClass(v),
		})
	}

	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label(),
			Class:       s.graph.GetEdgeClass(e),
			Source:      s.graph.EdgeSourceLabel(e),
			Destination: s.graph.EdgeDestinationLabel(e),
		})
	}
	return sub
}

func (s *Service) Neighbors(label string) SubGraph {

	principal := s.graph.GetVertexByLabel(label)
	if principal == nil {
		panic("vertex not found")
	}

	sub := SubGraph{
		Title: "Vizinhos de " + label,
		Principal: Vertex{
			Label: principal.Label(),
			Class: s.graph.GetVertexClass(principal),
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	rs := s.graph.Neighbors(principal)

	for _, v := range rs.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label(),
			Class: s.graph.GetVertexClass(v),
		})
	}

	for _, e := range rs.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label(),
			Class:       s.graph.GetEdgeClass(e),
			Source:      s.graph.EdgeSourceLabel(e),
			Destination: s.graph.EdgeDestinationLabel(e),
		})
	}

	return sub
}

func (s *Service) Path(label, destination string) SubGraph {
	res := s.graph.Path(label, destination)

	sub := SubGraph{
		Title: "Caminho de " + label + " para " + destination,
		Principal: Vertex{
			Label: res.Principal.Label(),
			Class: s.graph.GetVertexClass(res.Principal),
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
	}

	for _, v := range res.Vertices {
		sub.Vertices = append(sub.Vertices, Vertex{
			Label: v.Label(),
			Class: s.graph.GetVertexClass(v),
		})
	}
	for _, e := range res.Edges {
		sub.Edges = append(sub.Edges, Edge{
			Label:       e.Label(),
			Class:       s.graph.GetEdgeClass(e),
			Source:      s.graph.EdgeSourceLabel(e),
			Destination: s.graph.EdgeDestinationLabel(e),
		})
	}
	return sub
}
