package agent

import (
	"context"

	"github.com/mariomac/pipes/pkg/graph"
	"github.com/mariomac/pipes/pkg/node"

	"github.com/grafana/beyla/pkg/internal/netolly/ebpf"
	"github.com/grafana/beyla/pkg/internal/netolly/export"
	"github.com/grafana/beyla/pkg/internal/netolly/flow"
	"github.com/grafana/beyla/pkg/internal/netolly/transform/k8s"
)

// FlowsPipeline defines the different nodes in the Beyla's NetO11y module,
// as well as how they are interconnected
// TODO: add flow_printer node
type FlowsPipeline struct {
	MapTracer       `sendTo:"Deduper"`
	RingBufTracer   `sendTo:"Accounter"`
	Accounter       `sendTo:"Deduper"`
	Deduper         flow.Deduper `forwardTo:"CapacityLimiter"`
	CapacityLimiter `sendTo:"Decorator"`
	Decorator       `sendTo:"Kubernetes"`

	Kubernetes k8s.NetworkTransformConfig `sendTo:"Exporter"`

	Exporter export.MetricsConfig
}

type MapTracer struct{}
type RingBufTracer struct{}
type Accounter struct{}
type CapacityLimiter struct{}
type Decorator struct{}

// buildAndStartPipeline creates the ETL flow processing graph.
// For a more visual view, check the docs/architecture.md document.
func (f *Flows) buildAndStartPipeline(ctx context.Context) (graph.Graph, error) {

	alog := alog()
	alog.Debug("registering interfaces' listener in background")
	err := f.interfacesManager(ctx)
	if err != nil {
		return graph.Graph{}, err
	}

	alog.Debug("creating flows' processing graph")
	gb := graph.NewBuilder(node.ChannelBufferLen(f.cfg.ChannelBufferLen))

	graph.RegisterStart(gb, func(_ MapTracer) (node.StartFunc[[]*ebpf.Record], error) {
		return f.mapTracer.TraceLoop(ctx), nil
	})
	graph.RegisterStart(gb, func(_ RingBufTracer) (node.StartFunc[*ebpf.NetFlowRecordT], error) {
		return f.rbTracer.TraceLoop(ctx), nil
	})
	graph.RegisterMiddle(gb, func(_ Accounter) (node.MiddleFunc[*ebpf.NetFlowRecordT, []*ebpf.Record], error) {
		return f.accounter.Account, nil
	})
	graph.RegisterMiddle(gb, flow.DeduperProvider)
	graph.RegisterMiddle(gb, func(_ CapacityLimiter) (node.MiddleFunc[[]*ebpf.Record, []*ebpf.Record], error) {
		return (&flow.CapacityLimiter{}).Limit, nil
	})
	graph.RegisterMiddle(gb, func(_ Decorator) (node.MiddleFunc[[]*ebpf.Record, []*ebpf.Record], error) {
		return flow.Decorate(f.agentIP, f.interfaceNamer), nil
	})
	graph.RegisterMiddle(gb, k8s.NetworkTransform)

	graph.RegisterTerminal(gb, export.MetricsExporterProvider)

	var deduperExpireTime = f.cfg.NetworkFlows.DeduperFCExpiry
	if deduperExpireTime <= 0 {
		deduperExpireTime = 2 * f.cfg.NetworkFlows.CacheActiveTimeout
	}
	return gb.Build(&FlowsPipeline{
		Deduper: flow.Deduper{
			Type:       f.cfg.NetworkFlows.Deduper,
			ExpireTime: deduperExpireTime,
			JustMark:   f.cfg.NetworkFlows.DeduperJustMark,
		},
		Kubernetes: k8s.NetworkTransformConfig{Kubernetes: &f.cfg.Attributes.Kubernetes},
		// TODO: allow prometheus exporting
		Exporter: export.MetricsConfig{Metrics: &f.cfg.Metrics},
	})
}
