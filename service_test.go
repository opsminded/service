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
	BaseEdges: []graphlib.Edge{
		{
			Label:       "AB",
			Source:      graphlib.Vertex{Label: "A", Healthy: true},
			Destination: graphlib.Vertex{Label: "B", Healthy: true},
		},
		{
			Label:       "BC",
			Source:      graphlib.Vertex{Label: "B", Healthy: true},
			Destination: graphlib.Vertex{Label: "C", Healthy: true},
		},
	},
	BaseVertices: []graphlib.Vertex{
		{
			Label:   "A",
			Healthy: true,
		},
		{
			Label:   "B",
			Healthy: true,
		},
		{
			Label:   "C",
			Healthy: true,
		},
	},
}

func TestServiceBasics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	v, err := s.GetVertex(simpleEx.BaseVertices[0].Label)
	if err != nil {
		t.Fatal(err)
	}

	if v.Label != simpleEx.BaseVertices[0].Label {
		t.Fatal("Error")
	}
	cancel()
}

func TestSummary(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	sum := s.Summary()
	if sum.TotalVertices != len(simpleEx.BaseVertices) {
		t.Fatal("num vertices error")
	}

	if sum.TotalEdges != len(simpleEx.BaseEdges) {
		t.Fatal("num edges error")
	}

	if len(sum.UnhealthyVertices) != 0 {
		t.Fatal("num of unhealth error", len(sum.UnhealthyVertices))
	}
	cancel()
}

func TestNeighbors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	time.Sleep(2 * time.Second)

	nei, err := s.Neighbors("B")
	if err != nil {
		t.Fatal(err)
	}

	if len(nei.SubGraph.Vertices) != 3 {
		t.Fatal("2 nodes expected")
	}

	if len(nei.SubGraph.Edges) != 2 {
		t.Fatal("2 edges expected")
	}
	cancel()
}
