// This file was automatically generated by lister-gen

package machine

import (
	machine "github.com/gardener/node-controller-manager/pkg/apis/machine"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MachineSetLister helps list MachineSets.
type MachineSetLister interface {
	// List lists all MachineSets in the indexer.
	List(selector labels.Selector) (ret []*machine.MachineSet, err error)
	// Get retrieves the MachineSet from the index for a given name.
	Get(name string) (*machine.MachineSet, error)
	MachineSetListerExpansion
}

// machineSetLister implements the MachineSetLister interface.
type machineSetLister struct {
	indexer cache.Indexer
}

// NewMachineSetLister returns a new MachineSetLister.
func NewMachineSetLister(indexer cache.Indexer) MachineSetLister {
	return &machineSetLister{indexer: indexer}
}

// List lists all MachineSets in the indexer.
func (s *machineSetLister) List(selector labels.Selector) (ret []*machine.MachineSet, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*machine.MachineSet))
	})
	return ret, err
}

// Get retrieves the MachineSet from the index for a given name.
func (s *machineSetLister) Get(name string) (*machine.MachineSet, error) {
	key := &machine.MachineSet{ObjectMeta: v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(machine.Resource("machineset"), name)
	}
	return obj.(*machine.MachineSet), nil
}
