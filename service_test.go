package service_test

import (
	"testing"
	"time"

	"github.com/opsminded/service"
)

var simpleEx = &service.TestableExtractor{
	FrequencyDuration: time.Second,
	Edges: []service.Edge{
		{
			Label:       "AB",
			Source:      "A",
			Destination: "B",
		},
		{
			Label:       "BC",
			Source:      "B",
			Destination: "C",
		},
	},
	Vertices: []service.Vertex{
		{
			Label: "A",
		},
		{
			Label: "B",
		},
		{
			Label: "C",
		},
	},
}

func TestServiceBasics(t *testing.T) {
	s := service.New([]service.Extractor{simpleEx})
	s.Extract()

	v, err := s.GetVertex(simpleEx.Vertices[0].Label)
	if err != nil {
		t.Fatal(err)
	}

	if v.Label != simpleEx.Vertices[0].Label {
		t.Fatal("Error")
	}

	if _, err := s.GetVertex("x"); err == nil {
		t.Fatal("error expected")
	}
}

func TestSummary(t *testing.T) {
	s := service.New([]service.Extractor{simpleEx})
	s.Extract()

	sum := s.Summary()
	if sum.TotalVertex != len(simpleEx.Vertices) {
		t.Fatal("num vertices error")
	}

	if sum.TotalEdges != len(simpleEx.Edges) {
		t.Fatal("num edges error")
	}

	if len(sum.UnhealthVertex) != len(simpleEx.Vertices) {
		t.Fatal("num of unhealth error")
	}
}

func TestNeighbors(t *testing.T) {
	s := service.New([]service.Extractor{simpleEx})
	s.Extract()

	nei := s.Neighbors("B")
	if len(nei.Vertices) != 2 {
		t.Fatal("2 nodes expected")
	}

	if len(nei.Edges) != 2 {
		t.Fatal("2 edges expected")
	}
}
