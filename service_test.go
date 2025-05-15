package service_test

import (
	"testing"
	"time"

	"github.com/opsminded/graphlib/v2"
	"github.com/opsminded/service"
)

var simpleEx = &service.TestableExtractor{
	FrequencyDuration: time.Second,
	BaseEdges: []graphlib.Edge{
		{
			Source: "A",
			Target: "B",
		},
		{
			Source: "B",
			Target: "C",
		},
	},
	BaseVertices: []graphlib.Vertex{
		{
			Key:       "A",
			Label:     "A",
			Healthy:   true,
			LastCheck: time.Now().UnixNano(),
		},
		{
			Key:       "B",
			Label:     "B",
			Healthy:   true,
			LastCheck: time.Now().UnixNano(),
		},
		{
			Key:       "C",
			Label:     "C",
			Healthy:   true,
			LastCheck: time.Now().UnixNano(),
		},
	},
}

func TestServiceBasics(t *testing.T) {
	graph := graphlib.NewGraph()
	s := service.New(graph, []service.Extractor{simpleEx})

	v, err := s.GetVertex(simpleEx.BaseVertices[0].Label)
	if err != nil {
		t.Fatal(err)
	}

	if v.Label != simpleEx.BaseVertices[0].Label {
		t.Fatal("Error")
	}
}

// func TestSummary(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
// 	time.Sleep(2 * time.Second)

// 	sum := s.Summary()
// 	if sum.TotalVertices != len(simpleEx.BaseVertices) {
// 		t.Fatal("num vertices error")
// 	}

// 	if sum.TotalEdges != len(simpleEx.BaseEdges) {
// 		t.Fatal("num edges error")
// 	}

// 	if len(sum.UnhealthyVertices) != 0 {
// 		t.Fatal("num of unhealth error", len(sum.UnhealthyVertices))
// 	}
// 	cancel()
// }

func TestNeighbors(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// s := service.New(ctx, time.Second, []service.Extractor{simpleEx}, nil)
	// time.Sleep(2 * time.Second)

	// nei, err := s.Neighbors("B")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if len(nei.SubGraph.Vertices) != 3 {
	// 	t.Fatal("2 nodes expected")
	// }

	// if len(nei.SubGraph.Edges) != 2 {
	// 	t.Fatal("2 edges expected")
	// }
	// cancel()
}
