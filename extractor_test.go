package service_test

import (
	"sync"
	"time"

	"github.com/opsminded/service"
)

type TestableExtractor struct {
	FrequencyDuration time.Duration

	Edges    []service.Edge
	Vertices []service.Vertex

	edgePointer   int
	vertexPointer int

	mu sync.Mutex
}

var _ service.Extractor = (*TestableExtractor)(nil)

func (e *TestableExtractor) Frequency() time.Duration {
	return e.FrequencyDuration
}

func (e *TestableExtractor) NextEdge() service.Edge {
	e.mu.Lock()
	defer e.mu.Unlock()
	edge := e.Edges[e.edgePointer]
	e.edgePointer++
	return edge
}

func (e *TestableExtractor) HasNextEdge() bool {
	return e.edgePointer < len(e.Edges)
}

func (e *TestableExtractor) NextVertex() service.Vertex {
	e.mu.Lock()
	defer e.mu.Unlock()
	vertex := e.Vertices[e.vertexPointer]
	e.vertexPointer++
	return vertex
}

func (e *TestableExtractor) HasNextVertex() bool {
	return e.vertexPointer < len(e.Vertices)
}

func (e *TestableExtractor) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.edgePointer = 0
	e.vertexPointer = 0
}
