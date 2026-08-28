package elements

import (
	"laatoo.io/sdk/server/components/knowledge"
	"laatoo.io/sdk/server/core"
)

// KnowledgeManager is the server element that resolves a knowledge graph by alias, obtained with
// ctx.GetServerElement(core.ServerElementKnowledgeManager).(elements.KnowledgeManager).
//
// It mirrors AgentManager's memory sub-registry: a GraphType-keyed provider registry that lets a
// storage backend register itself instead of a consumer plugin switching on a hardcoded backend
// name, plus an alias-keyed registry of the graphs actually built. Both registries are populated
// at startup by module loading. The request-time surface -- RunSPARQL and Reload -- lives on the
// resolved knowledge.Graph itself, not here.
type KnowledgeManager interface {
	core.ServerElement

	// RegisterGraphProvider claims a GraphType for a storage backend. Call it from the backend
	// plugin's Initialize, before any graph naming that type is built.
	RegisterGraphProvider(ctx core.ServerContext, graphType knowledge.GraphType, p knowledge.GraphProvider) error

	// GetGraphProvider returns the provider registered for graphType, or an error if none is --
	// the lookup RegisterGraphProvider exists to make possible.
	GetGraphProvider(ctx core.ServerContext, graphType knowledge.GraphType) (knowledge.GraphProvider, error)

	// RegisterGraph registers an already-built graph under alias, mirroring AgentManager's direct
	// RegisterAgent path for pro-code components that build their own instance rather than going
	// through a type-keyed provider.
	RegisterGraph(ctx core.ServerContext, alias string, g knowledge.Graph) error

	// GetGraph returns the graph registered under alias, or an error if none is.
	GetGraph(ctx core.ServerContext, alias string) (knowledge.Graph, error)

	// ListGraphs returns every registered graph keyed by alias.
	ListGraphs(ctx core.ServerContext) map[string]knowledge.Graph
}
