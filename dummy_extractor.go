package service

import (
	"sync"
	"time"

	"github.com/opsminded/graphlib/v2"
)

type TestableExtractor struct {
	FrequencyDuration time.Duration

	BaseEdges    []graphlib.Edge
	BaseVertices []graphlib.Vertex

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
	edge := e.BaseEdges[e.edgePointer]
	e.edgePointer++
	return edge
}

func (e *TestableExtractor) HasNextEdge() bool {
	return e.edgePointer < len(e.BaseEdges)
}

func (e *TestableExtractor) NextVertex() graphlib.Vertex {
	e.mu.Lock()
	defer e.mu.Unlock()
	vertex := e.BaseVertices[e.vertexPointer]
	e.vertexPointer++
	return vertex
}

func (e *TestableExtractor) HasNextVertex() bool {
	return e.vertexPointer < len(e.BaseVertices)
}

func (e *TestableExtractor) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.edgePointer = 0
	e.vertexPointer = 0
}
