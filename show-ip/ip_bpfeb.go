// Code generated by bpf2go; DO NOT EDIT.
//go:build mips || mips64 || ppc64 || s390x

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type ipInfo struct {
	SourceIp uint32
	DestIp   uint32
	Ttl      uint8
	Protocol uint8
	_        [2]byte
}

// loadIp returns the embedded CollectionSpec for ip.
func loadIp() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_IpBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load ip: %w", err)
	}

	return spec, err
}

// loadIpObjects loads ip and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*ipObjects
//	*ipPrograms
//	*ipMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadIpObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadIp()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// ipSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipSpecs struct {
	ipProgramSpecs
	ipMapSpecs
}

// ipSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipProgramSpecs struct {
	GetIps *ebpf.ProgramSpec `ebpf:"get_ips"`
}

// ipMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipMapSpecs struct {
	Ips *ebpf.MapSpec `ebpf:"ips"`
}

// ipObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadIpObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipObjects struct {
	ipPrograms
	ipMaps
}

func (o *ipObjects) Close() error {
	return _IpClose(
		&o.ipPrograms,
		&o.ipMaps,
	)
}

// ipMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadIpObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipMaps struct {
	Ips *ebpf.Map `ebpf:"ips"`
}

func (m *ipMaps) Close() error {
	return _IpClose(
		m.Ips,
	)
}

// ipPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadIpObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipPrograms struct {
	GetIps *ebpf.Program `ebpf:"get_ips"`
}

func (p *ipPrograms) Close() error {
	return _IpClose(
		p.GetIps,
	)
}

func _IpClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed ip_bpfeb.o
var _IpBytes []byte
