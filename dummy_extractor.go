package service

import (
	"sync"
	"time"

	"github.com/opsminded/graphlib"
)

type TestableExtractor struct {
	FrequencyDuration time.Duration

	Edges    []graphlib.Edge
	Vertices []graphlib.Vertex

	edgePointer   int
	vertexPointer int

	mu sync.Mutex
}

var _ Extractor = (*TestableExtractor)(nil)

func (e *TestableExtractor) Frequency() time.Duration {
	return e.FrequencyDuration
}

func (e *TestableExtractor) NextEdge() graphlib.Edge {
	e.mu.Lock()
	defer e.mu.Unlock()
	edge := e.Edges[e.edgePointer]
	e.edgePointer++
	return edge
}

func (e *TestableExtractor) HasNextEdge() bool {
	return e.edgePointer < len(e.Edges)
}

func (e *TestableExtractor) NextVertex() graphlib.Vertex {
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
