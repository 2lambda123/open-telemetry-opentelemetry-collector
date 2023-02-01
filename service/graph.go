// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service // import "go.opentelemetry.io/collector/service"

import (
	"context"
	"net/http"

	"go.uber.org/multierr"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"go.opentelemetry.io/collector/component"
)

var _ pipelines = (*pipelinesGraph)(nil)

type pipelinesGraph struct {
	// All component instances represented as nodes, with directed edges indicating data flow.
	componentGraph *simple.DirectedGraph

	// Keep track of how nodes relate to pipelines, so we can declare edges in the graph.
	pipelines map[component.ID]*pipelineNodes
}

func buildPipelinesGraph(ctx context.Context, set pipelinesSettings) (pipelines, error) {
	pipelines := &pipelinesGraph{
		componentGraph: simple.NewDirectedGraph(),
		pipelines:      make(map[component.ID]*pipelineNodes, len(set.PipelineConfigs)),
	}
	for pipelineID := range set.PipelineConfigs {
		pipelines.pipelines[pipelineID] = &pipelineNodes{
			receivers: make(map[int64]graph.Node),
			exporters: make(map[int64]graph.Node),
		}
	}
	pipelines.createNodes(set)
	pipelines.createEdges()
	return pipelines, pipelines.buildComponents(ctx, set)
}

// Creates a node for each instance of a component and adds it to the graph
func (g *pipelinesGraph) createNodes(set pipelinesSettings) {
	// Keep track of connectors and where they are used. (map[connectorID][]pipelineID)
	connectorsAsExporter := make(map[component.ID][]component.ID)
	connectorsAsReceiver := make(map[component.ID][]component.ID)

	for pipelineID, pipelineCfg := range set.PipelineConfigs {
		pipe := g.pipelines[pipelineID]
		for _, recvID := range pipelineCfg.Receivers {
			if set.Connectors.IsConfigured(recvID) {
				connectorsAsReceiver[recvID] = append(connectorsAsReceiver[recvID], pipelineID)
				continue
			}
			rcvrNode := g.createReceiver(pipelineID, recvID)
			pipe.receivers[rcvrNode.ID()] = rcvrNode
		}

		pipe.capabilitiesNode = newCapabilitiesNode(pipelineID)

		for _, procID := range pipelineCfg.Processors {
			pipe.processors = append(pipe.processors, g.createProcessor(pipelineID, procID))
		}

		pipe.fanOutNode = newFanOutNode(pipelineID)

		for _, exprID := range pipelineCfg.Exporters {
			if set.Connectors.IsConfigured(exprID) {
				connectorsAsExporter[exprID] = append(connectorsAsExporter[exprID], pipelineID)
				continue
			}
			expNode := g.createExporter(pipelineID, exprID)
			pipe.exporters[expNode.ID()] = expNode
		}
	}

	for connID, exprPipelineIDs := range connectorsAsExporter {
		for _, eID := range exprPipelineIDs {
			for _, rID := range connectorsAsReceiver[connID] {
				connNode := g.createConnector(eID, rID, connID)
				g.pipelines[eID].exporters[connNode.ID()] = connNode
				g.pipelines[rID].receivers[connNode.ID()] = connNode
			}
		}
	}
}

func (g *pipelinesGraph) createReceiver(pipelineID, recvID component.ID) *receiverNode {
	rcvrNode := newReceiverNode(pipelineID, recvID)
	if node := g.componentGraph.Node(rcvrNode.ID()); node != nil {
		return node.(*receiverNode)
	}
	g.componentGraph.AddNode(rcvrNode)
	return rcvrNode
}

func (g *pipelinesGraph) createProcessor(pipelineID, procID component.ID) *processorNode {
	procNode := newProcessorNode(pipelineID, procID)
	g.componentGraph.AddNode(procNode)
	return procNode
}

func (g *pipelinesGraph) createExporter(pipelineID, exprID component.ID) *exporterNode {
	expNode := newExporterNode(pipelineID, exprID)
	if node := g.componentGraph.Node(expNode.ID()); node != nil {
		return node.(*exporterNode)
	}
	g.componentGraph.AddNode(expNode)
	return expNode
}

func (g *pipelinesGraph) createConnector(exprPipelineID, rcvrPipelineID, connID component.ID) *connectorNode {
	connNode := newConnectorNode(exprPipelineID.Type(), rcvrPipelineID.Type(), connID)
	if node := g.componentGraph.Node(connNode.ID()); node != nil {
		return node.(*connectorNode)
	}
	g.componentGraph.AddNode(connNode)
	return connNode
}

func (g *pipelinesGraph) createEdges() {
	for _, pg := range g.pipelines {
		for _, receiver := range pg.receivers {
			g.componentGraph.SetEdge(g.componentGraph.NewEdge(receiver, pg.capabilitiesNode))
		}

		var from, to graph.Node
		from = pg.capabilitiesNode
		for _, processor := range pg.processors {
			to = processor
			g.componentGraph.SetEdge(g.componentGraph.NewEdge(from, to))
			from = processor
		}
		to = pg.fanOutNode
		g.componentGraph.SetEdge(g.componentGraph.NewEdge(from, to))

		for _, exporter := range pg.exporters {
			g.componentGraph.SetEdge(g.componentGraph.NewEdge(pg.fanOutNode, exporter))
		}
	}
}

func (g *pipelinesGraph) buildComponents(ctx context.Context, set pipelinesSettings) error {
	nodes, err := topo.Sort(g.componentGraph)
	if err != nil {
		return err // TODO clean up error message
	}

	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		switch n := node.(type) {
		case *receiverNode:
			err = n.build(ctx, set.Telemetry, set.BuildInfo, set.Receivers, g.nextConsumers(n.ID()))
		case *processorNode:
			err = n.build(ctx, set.Telemetry, set.BuildInfo, set.Processors, g.nextConsumers(n.ID())[0])
		case *connectorNode:
			err = n.build(ctx, set.Telemetry, set.BuildInfo, set.Connectors, g.nextConsumers(n.ID()))
		case *exporterNode:
			err = n.build(ctx, set.Telemetry, set.BuildInfo, set.Exporters)
		case *capabilitiesNode:
			n.build(g.nextConsumers(n.ID())[0], g.pipelines[n.pipelineID].processors)
		case *fanOutNode:
			n.build(g.nextConsumers(n.ID()))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Find all nodes
func (g *pipelinesGraph) nextConsumers(nodeID int64) []baseConsumer {
	nextNodes := g.componentGraph.From(nodeID)
	nexts := make([]baseConsumer, 0, nextNodes.Len())
	for nextNodes.Next() {
		nexts = append(nexts, nextNodes.Node().(consumerNode).getConsumer())
	}
	return nexts
}

// A node-based representation of a pipeline configuration.
type pipelineNodes struct {
	// Use map to assist with deduplication of connector instances.
	receivers map[int64]graph.Node

	// The node to which receivers emit. Passes through to processors.
	// Easily accessible as the first node in a pipeline.
	*capabilitiesNode

	// The order of processors is very important. Therefore use a slice for processors.
	processors []*processorNode

	// Emits to exporters.
	*fanOutNode

	// Use map to assist with deduplication of connector instances.
	exporters map[int64]graph.Node
}

func (g *pipelinesGraph) StartAll(ctx context.Context, host component.Host) error {
	nodes, err := topo.Sort(g.componentGraph)
	if err != nil {
		return err
	}

	// Start in reverse topological order so that downstream components
	// are started before upstream components. This ensures that each
	// component's consumer is ready to consume.
	for i := len(nodes) - 1; i >= 0; i-- {
		comp, ok := nodes[i].(component.Component)
		if !ok {
			// Skip capabilities/fanout nodes
			continue
		}
		if compErr := comp.Start(ctx, host); compErr != nil {
			return compErr
		}
	}
	return nil
}

func (g *pipelinesGraph) ShutdownAll(ctx context.Context) error {
	nodes, err := topo.Sort(g.componentGraph)
	if err != nil {
		return err
	}

	// Stop in topological order so that upstream components
	// are stopped before downstream components.  This ensures
	// that each component has a chance to drain to it's consumer
	// before the consumer is stopped.
	var errs error
	for i := 0; i < len(nodes); i++ {
		comp, ok := nodes[i].(component.Component)
		if !ok {
			// Skip capabilities/fanout nodes
			continue
		}
		errs = multierr.Append(errs, comp.Shutdown(ctx))
	}
	return errs
}

func (g *pipelinesGraph) GetExporters() map[component.DataType]map[component.ID]component.Component {
	// TODO actual implementation
	return make(map[component.DataType]map[component.ID]component.Component)
}

func (g *pipelinesGraph) HandleZPages(w http.ResponseWriter, r *http.Request) {
	// TODO actual implementation
}
