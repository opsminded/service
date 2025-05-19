package service

import (
	"log/slog"

	"github.com/opsminded/graphlib/v2"
)

type VertexAttribute struct {
	Type        string
	Description string
	Value       any
}

type QueryResult struct {
	Title     string
	Principal graphlib.Vertex
	SubGraph  graphlib.Subgraph
}

type Summary struct {
	TotalEdges             int
	TotalVertices          int
	TotalHealthyVertices   int
	TotalUnhealthyVertices int
	UnhealthyVertices      []graphlib.Vertex
}

type Service struct {
	graph *graphlib.Graph
}

func New(graph *graphlib.Graph) *Service {
	service := &Service{
		graph: graph,
	}
	return service
}

func (s *Service) GetVertex(key string) (graphlib.Vertex, error) {
	slog.Debug("service.GetVertex", slog.String("key", key))

	v, err := s.graph.GetVertex(key)
	if err != nil {
		slog.Error("service.GetVertex", slog.String("key", key), slog.String("error", err.Error()))
		return graphlib.Vertex{}, err
	}
	return v, nil
}

func (s *Service) GetVertexAttributes(key string) ([]VertexAttribute, error) {
	slog.Debug("service.GetVertexAttributes", slog.String("key", key))

	_, err := s.graph.GetVertex(key)
	if err != nil {
		slog.Error("service.GetVertexAttributes", slog.String("key", key), slog.String("error", err.Error()))
		return nil, err
	}

	return []VertexAttribute{
		{Type: "string", Description: "Grupo de suporte", Value: "234 - DITI"},
		{Type: "string", Description: "Gerente", Value: "João Paulo"},
		{Type: "string", Description: "Cluster", Value: "JAHSGDSALHD"},
		{Type: "string", Description: "Ambiente", Value: "Produção"},
		{Type: "link", Description: "Link monitoração", Value: "http://google.com.br"},
		{Type: "link", Description: "Repositório Github", Value: "http://github.com"},
		{Type: "link", Description: "Consultar incidentes", Value: "http://google.com.br"},
	}, nil
}

func (s *Service) SetVertexHealth(key string, health bool) error {
	slog.Debug("service.SetVertexHealth", slog.String("key", key), slog.Bool("health", health))

	err := s.graph.SetVertexHealth(key, health)
	if err != nil {
		slog.Error("service.SetVertexHealth", slog.String("key", key), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Service) ClearGraphHealthyStatus() {
	slog.Debug("service.ClearGraphHealthyStatus")
	s.graph.ClearGraphHealthyStatus()
}

func (s *Service) Summary() Summary {
	slog.Debug("service.Summary")

	stats := s.graph.GraphStats()

	sum := Summary{
		TotalEdges:    stats.TotalEdges,
		TotalVertices: stats.TotalVertices,

		TotalHealthyVertices:   stats.TotalHealthyVertices,
		TotalUnhealthyVertices: stats.TotalUnhealthyVertices,

		UnhealthyVertices: stats.UnhealthyVertices,
	}
	return sum
}

func (s *Service) VertexDependencies(key string, all bool) (QueryResult, error) {
	slog.Debug("service.VertexDependencies", slog.String("key", key), slog.Bool("all", all))

	p, err := s.graph.GetVertex(key)
	if err != nil {
		slog.Error("service.VertexDependencies, graph.GetVertex", slog.String("key", key), slog.Bool("all", all), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	dep, err := s.graph.VertexDependencies(key, all)
	if err != nil {
		slog.Error("service.VertexDependencies, graph.VertexDependencies", slog.String("key", key), slog.Bool("all", all), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Dependências de " + p.Label,
		Principal: p,
		SubGraph:  dep,
	}
	return sub, nil
}

func (s *Service) VertexDependents(key string, all bool) (QueryResult, error) {
	slog.Debug("service.VertexDependents", slog.String("key", key), slog.Bool("all", all))

	p, err := s.graph.GetVertex(key)
	if err != nil {
		slog.Error("service.VertexDependents, graph.GetVertex", slog.String("key", key), slog.Bool("all", all), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	dep, err := s.graph.VertexDependents(key, all)
	if err != nil {
		slog.Error("service.VertexDependents, graph.VertexDependents", slog.String("key", key), slog.Bool("all", all), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Dependentes de " + p.Label,
		Principal: p,
		SubGraph:  dep,
	}
	return sub, nil
}

func (s *Service) VertexNeighbors(key string) (QueryResult, error) {
	slog.Debug("service.VertexNeighbors", slog.String("key", key))

	p, err := s.graph.GetVertex(key)
	if err != nil {
		slog.Error("service.VertexNeighbors, graph.GetVertex", slog.String("key", key))
		return QueryResult{}, err
	}

	neighbors, err := s.graph.VertexNeighbors(key)
	if err != nil {
		slog.Error("service.VertexNeighbors, graph.VertexNeighbors", slog.String("key", key), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Vizinhos de " + p.Label,
		Principal: p,
		SubGraph:  neighbors,
	}
	return sub, nil
}

func (s *Service) Path(kSrc, ktgt string) (QueryResult, error) {
	slog.Debug("service.Path", slog.String("kSrc", kSrc), slog.String("ktgt", ktgt))

	src, err := s.graph.GetVertex(kSrc)
	if err != nil {
		slog.Error("service.Path, graph.GetVertex", slog.String("kSrc", kSrc), slog.String("ktgt", ktgt), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	tgt, err := s.graph.GetVertex(ktgt)
	if err != nil {
		slog.Error("service.Path, graph.GetVertex", slog.String("kSrc", kSrc), slog.String("ktgt", ktgt), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	path, err := s.graph.Path(kSrc, ktgt)
	if err != nil {
		slog.Error("service.Path, graph.Path", slog.String("kSrc", kSrc), slog.String("ktgt", ktgt), slog.String("error", err.Error()))
		return QueryResult{}, err
	}

	sub := QueryResult{
		Title:     "Caminhos de " + src.Label + " para " + tgt.Label,
		Principal: src,
		SubGraph:  path,
	}
	return sub, nil
}
