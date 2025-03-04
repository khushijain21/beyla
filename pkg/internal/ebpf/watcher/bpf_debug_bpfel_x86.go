// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package watcher

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type bpf_debugHttpConnectionMetadataT struct {
	Pid struct {
		HostPid   uint32
		UserPid   uint32
		Namespace uint32
	}
	Type uint8
}

type bpf_debugPidConnectionInfoT struct {
	Conn struct {
		S_addr [16]uint8
		D_addr [16]uint8
		S_port uint16
		D_port uint16
	}
	Pid uint32
}

type bpf_debugPidKeyT struct {
	Pid       uint32
	Namespace uint32
}

type bpf_debugWatchInfoT struct {
	Flags   uint64
	Payload uint64
}

// loadBpf_debug returns the embedded CollectionSpec for bpf_debug.
func loadBpf_debug() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Bpf_debugBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf_debug: %w", err)
	}

	return spec, err
}

// loadBpf_debugObjects loads bpf_debug and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*bpf_debugObjects
//	*bpf_debugPrograms
//	*bpf_debugMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpf_debugObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf_debug()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpf_debugSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugSpecs struct {
	bpf_debugProgramSpecs
	bpf_debugMapSpecs
}

// bpf_debugSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugProgramSpecs struct {
	KprobeSysBind *ebpf.ProgramSpec `ebpf:"kprobe_sys_bind"`
}

// bpf_debugMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugMapSpecs struct {
	FilteredConnections *ebpf.MapSpec `ebpf:"filtered_connections"`
	PidCache            *ebpf.MapSpec `ebpf:"pid_cache"`
	ValidPids           *ebpf.MapSpec `ebpf:"valid_pids"`
	WatchEvents         *ebpf.MapSpec `ebpf:"watch_events"`
}

// bpf_debugObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugObjects struct {
	bpf_debugPrograms
	bpf_debugMaps
}

func (o *bpf_debugObjects) Close() error {
	return _Bpf_debugClose(
		&o.bpf_debugPrograms,
		&o.bpf_debugMaps,
	)
}

// bpf_debugMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugMaps struct {
	FilteredConnections *ebpf.Map `ebpf:"filtered_connections"`
	PidCache            *ebpf.Map `ebpf:"pid_cache"`
	ValidPids           *ebpf.Map `ebpf:"valid_pids"`
	WatchEvents         *ebpf.Map `ebpf:"watch_events"`
}

func (m *bpf_debugMaps) Close() error {
	return _Bpf_debugClose(
		m.FilteredConnections,
		m.PidCache,
		m.ValidPids,
		m.WatchEvents,
	)
}

// bpf_debugPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugPrograms struct {
	KprobeSysBind *ebpf.Program `ebpf:"kprobe_sys_bind"`
}

func (p *bpf_debugPrograms) Close() error {
	return _Bpf_debugClose(
		p.KprobeSysBind,
	)
}

func _Bpf_debugClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_debug_bpfel_x86.o
var _Bpf_debugBytes []byte
