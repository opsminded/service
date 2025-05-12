package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/opsminded/graphlib"
	"github.com/opsminded/service"
)

var simpleEx = &service.TestableExtractor{
	FrequencyDuration: time.Second,
	Edges: []graphlib.Edge{
		{
			Label:       "AB",
			Source:      graphlib.Vertex{Label: "A"},
			Destination: graphlib.Vertex{Label: "B"},
		},
		{
			Label:       "BC",
			Source:      graphlib.Vertex{Label: "B"},
			Destination: graphlib.Vertex{Label: "C"},
		},
	},
	Vertices: []graphlib.Vertex{
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
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	v, err := s.GetVertex(simpleEx.Vertices[0].Label)
	if err != nil {
		t.Fatal(err)
	}

	if v.Label != simpleEx.Vertices[0].Label {
		t.Fatal("Error")
	}
	cancel()
}

func TestSummary(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	sum := s.Summary()
	if sum.TotalVertex != len(simpleEx.Vertices) {
		t.Fatal("num vertices error")
	}

	if sum.TotalEdges != len(simpleEx.Edges) {
		t.Fatal("num edges error")
	}

	if len(sum.UnhealthVertex) != 0 {
		t.Fatal("num of unhealth error")
	}
	cancel()
}

func TestNeighbors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	nei := s.Neighbors("B")
	if len(nei.SubGraph.Vertices) != 3 {
		t.Fatal("2 nodes expected")
	}

	if len(nei.SubGraph.Edges) != 2 {
		t.Fatal("2 edges expected")
	}
	cancel()
}
