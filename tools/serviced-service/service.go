// Copyright 2020 The Serviced Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/control-center/serviced/domain/service"
)

// ServiceList is a collection of Service entities
type ServiceList []service.Service

// NewServiceList returns a new ServiceList entity.
func NewServiceList() *ServiceList {
	return &ServiceList{}
}

// // NewServiceListFromJSON returns a new ServiceList entity populated by the given JSON string.
// func NewServiceListFromJSON(b []byte) (*ServiceList, error) {
// 	// var decoded []service.Service
// 	var list ServiceList
// 	if err := json.Unmarshal(b, &list); err != nil {
// 		return nil, err
// 	}
// 	return &list, nil
// }

// UnmarshalJSON transforms a JSON encoded byte array into a ServiceList
func (ss *ServiceList) UnmarshalJSON(b []byte) error {
	var decoded []service.Service
	if err := json.Unmarshal(b, &decoded); err != nil {
		return err
	}
	ss.Append(decoded...)
	return nil
}

// Append adds a service to the store.
func (ss *ServiceList) Append(items ...service.Service) {
	for _, svc := range items {
		ss.Put(svc)
	}
	// *ss = append(*ss, items...)
}

// Put adds the service.
// This function will overwrite a duplicate service.
func (ss *ServiceList) Put(svc service.Service) {
	var posn = -1
	for idx, existing := range *ss {
		if existing.ID == svc.ID {
			posn = idx
			break
		}
	}
	if posn >= 0 {
		(*ss)[posn] = svc
	} else {
		*ss = append(*ss, svc)
	}
}

// Len returns the number of items in the store.
func (ss *ServiceList) Len() int {
	return len(*ss)
}

// At returns the service at the given index.
// An error is returned if the position is out of range.
func (ss *ServiceList) At(posn uint) (*service.Service, error) {
	if posn >= uint(len(*ss)) {
		return nil, fmt.Errorf("Index out of range At(%v) with length %v", posn, len(*ss))
	}
	return &(*ss)[posn], nil
}

// Get returns the identified service.
// An error is returned if the service is not found.
func (ss *ServiceList) Get(serviceID string) (*service.Service, error) {
	for idx, svc := range *ss {
		if svc.ID == serviceID {
			return &(*ss)[idx], nil
		}
	}
	return nil, fmt.Errorf("Service not found with ID '%v'", serviceID)
}

// GetService returns a service having the given ID.
func (ss *ServiceList) GetService(serviceID string) (service.Service, error) {
	svc, err := ss.Get(serviceID)
	if err != nil {
		return service.Service{}, err
	}
	return *svc, nil
}

// GetServicePath returns the service's path.
func (ss *ServiceList) GetServicePath(serviceID string) (string, error) {
	var err error
	var svc *service.Service

	svc, err = ss.Get(serviceID)
	if err != nil {
		return "", err
	}
	if svc.ParentServiceID == "" {
		return fmt.Sprintf("/%s", svc.Name), nil
	}

	var segments = []string{svc.Name}
	for svc.ParentServiceID != "" {
		svc, err = ss.Get(svc.ParentServiceID)
		if err != nil {
			return "", err
		}
		segments = append([]string{svc.Name}, segments...)
	}
	return fmt.Sprintf("/%s", strings.Join(segments, "/")), nil
}

// GetTenantID returns the tenant ID.
// The tenant ID is the same ID as the service having no parent.
func (ss *ServiceList) GetTenantID() string {
	for _, svc := range *ss {
		if svc.ParentServiceID == "" {
			return svc.ID
		}
	}
	return ""
}

// FindChild returns the service on the indicated path
func (ss *ServiceList) FindChild(parentID, name string) (*service.Service, error) {
	parentID = strings.TrimSpace(parentID)
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, fmt.Errorf("Empty service name not allowed")
	}

	for _, svc := range *ss {
		if svc.ParentServiceID == parentID && svc.Name == name {
			return &svc, nil
		}
	}
	return nil, nil
}

// FindChildService returns the named service having the indicated parent service.
func (ss *ServiceList) FindChildService(parentID, name string) (service.Service, error) {
	svc, err := ss.FindChild(parentID, name)
	if err != nil {
		return service.Service{}, err
	}
	return *svc, nil
}

// Iterator returns an iterator over the ServiceList.
func (ss *ServiceList) Iterator() *ServiceListIterator {
	return &ServiceListIterator{
		store: ss,
		posn:  0,
	}
}

// ServiceListIterator is an iterator over a ServiceList.
type ServiceListIterator struct {
	store *ServiceList
	posn  uint
}

// Next advances the iterator and returns true if successful.
func (ssi *ServiceListIterator) Next() bool {
	if ssi.posn == uint(len(*(ssi.store))) {
		return false
	}
	ssi.posn = ssi.posn + 1
	return true
}

// Item returns the item the iterator currently points at.
func (ssi *ServiceListIterator) Item() *service.Service {
	value, err := ssi.store.At(ssi.posn - 1)
	if err != nil {
		panic(err)
	}
	return value
}
