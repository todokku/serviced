// Copyright 2017 The Serviced Authors.
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

// +build unit

package servicestatemanager

import (
	"fmt"
	"testing"

	"github.com/control-center/serviced/domain/service"

	. "gopkg.in/check.v1"
)

func TestServiceStateManager(t *testing.T) { TestingT(t) }

type ServiceStateManagerSuite struct {
	serviceStateManager ServiceStateManager
}

var _ = Suite(&ServiceStateManagerSuite{})

func (s *ServiceStateManagerSuite) SetUpSuite(c *C) {
	s.serviceStateManager = ServiceStateManager{
		tenantQueues: make(map[string]*ServiceStateQueue),
	}
}

func getTestServicesOne() []*service.Service {
	return []*service.Service{
		&service.Service{
			ID:                     "A",
			DesiredState:           1,
			EmergencyShutdownLevel: 0,
			StartLevel:             2,
		},
		&service.Service{
			ID:                     "B",
			DesiredState:           1,
			EmergencyShutdownLevel: 1,
			StartLevel:             3,
		},
		&service.Service{
			ID:                     "C",
			DesiredState:           1,
			EmergencyShutdownLevel: 2,
			StartLevel:             2,
		},
	}
}

func getTestServicesTwo() []*service.Service {
	return []*service.Service{
		&service.Service{
			ID:                     "D",
			DesiredState:           1,
			EmergencyShutdownLevel: 0,
			StartLevel:             2,
		},
		&service.Service{
			ID:                     "E",
			DesiredState:           1,
			EmergencyShutdownLevel: 1,
			StartLevel:             3,
		},
		&service.Service{
			ID:                     "F",
			DesiredState:           1,
			EmergencyShutdownLevel: 2,
			StartLevel:             2,
		},
	}
}

func getTestServicesThree() []*service.Service {
	return []*service.Service{
		&service.Service{
			ID:                     "A",
			DesiredState:           1,
			EmergencyShutdownLevel: 0,
			StartLevel:             2,
		},
		&service.Service{
			ID:                     "D",
			DesiredState:           1,
			EmergencyShutdownLevel: 0,
			StartLevel:             2,
		},
		&service.Service{
			ID:                     "G",
			DesiredState:           1,
			EmergencyShutdownLevel: 1,
			StartLevel:             3,
		},
		&service.Service{
			ID:                     "H",
			DesiredState:           1,
			EmergencyShutdownLevel: 2,
			StartLevel:             2,
		},
	}
}

func (s *ServiceStateManagerSuite) TestServiceStateManager_ScheduleServices_NoErr(c *C) {

	// Test that the batch has been added to the batchQueue
	// and split by nomral start level
	tenantID := "tenant"

	queue := &ServiceStateQueue{
		currentBatch: PendingServiceStateChangeBatch{},
	}
	s.serviceStateManager.tenantQueues[tenantID] = queue

	err := s.serviceStateManager.ScheduleServices(getTestServicesOne(), tenantID, service.SVCRun, false)
	if err != nil {
		c.Fatalf("Error in TestScheduleServices: %v\n", err)
	}

	c.Assert(s.serviceStateManager.tenantQueues[tenantID].batchQueue, DeepEquals, []ServiceStateChangeBatch{
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "A",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
				&service.Service{
					ID:                     "C",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "B",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
	})

	// Test that:
	// 1. The batch has been added to the batchQueue
	// 2. The batch was split by emergency shutdown level
	// 3. The emergency batches were moved to the front of the queue,
	// 4. The existing batches have been purged of the emergency-shutdown services
	err = s.serviceStateManager.ScheduleServices(getTestServicesOne(), tenantID, service.SVCStop, true)
	if err != nil {
		c.Fatalf("Error in TestScheduleServices: %v\n", err)
	}
	c.Logf("%s\n", s.serviceStateManager.tenantQueues[tenantID].batchQueue)
	c.Assert(s.serviceStateManager.tenantQueues[tenantID].batchQueue, DeepEquals, []ServiceStateChangeBatch{
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "B",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "C",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "A",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
	})

	// Test that trying to start a batch that has been scheduled for emergency shutdown has no effect on the queue
	err = s.serviceStateManager.ScheduleServices(getTestServicesOne(), tenantID, service.SVCRun, false)
	if err != nil {
		c.Fatalf("Error in TestScheduleServices: %v\n", err)
	}

	c.Assert(s.serviceStateManager.tenantQueues[tenantID].batchQueue, DeepEquals, []ServiceStateChangeBatch{
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "B",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "C",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "A",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
	})

	// Test that adding a non-conflicting non-emergency batch gets split by start level and appended to the queue
	err = s.serviceStateManager.ScheduleServices(getTestServicesTwo(), tenantID, service.SVCRun, false)
	if err != nil {
		c.Fatalf("Error in TestScheduleServices: %v\n", err)
	}

	c.Assert(s.serviceStateManager.tenantQueues[tenantID].batchQueue, DeepEquals, []ServiceStateChangeBatch{
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "B",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "C",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "A",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "D",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
				&service.Service{
					ID:                     "F",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "E",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
	})

	// Add a non-emergency batch with some conflicting and some non-conflicting services, and make sure that:
	//  1. The conflicting services are removed from the incoming batch
	//  2. The non-conflicting services are merged with the end of the queue based on start level
	err = s.serviceStateManager.ScheduleServices(getTestServicesThree(), tenantID, service.SVCRun, false)
	if err != nil {
		c.Fatalf("Error in TestScheduleServices: %v\n", err)
	}

	for _, batch := range s.serviceStateManager.tenantQueues[tenantID].batchQueue {
		s.LogBatch(c, batch)
	}

	c.Assert(s.serviceStateManager.tenantQueues[tenantID].batchQueue, DeepEquals, []ServiceStateChangeBatch{
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "B",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "C",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "A",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
			},
			desiredState: 0,
			emergency:    true,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "D",
					DesiredState:           1,
					EmergencyShutdownLevel: 0,
					StartLevel:             2,
				},
				&service.Service{
					ID:                     "F",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
				&service.Service{
					ID:                     "H",
					DesiredState:           1,
					EmergencyShutdownLevel: 2,
					StartLevel:             2,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
		ServiceStateChangeBatch{
			services: []*service.Service{
				&service.Service{
					ID:                     "E",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
				&service.Service{
					ID:                     "G",
					DesiredState:           1,
					EmergencyShutdownLevel: 1,
					StartLevel:             3,
				},
			},
			desiredState: 1,
			emergency:    false,
		},
	})

}

func (s *ServiceStateManagerSuite) CompareBatches(c *C, a, b ServiceStateChangeBatch) bool {
	sameVals := true
	if a.desiredState != b.desiredState {
		c.Logf("desiredState mismatch, a: %v b: %v", a.desiredState, b.desiredState)
		sameVals = false
	}
	if a.emergency != b.emergency {
		c.Logf("emergency mismatch, a: %v b: %v", a.emergency, b.emergency)
		sameVals = false
	}
	for n, svc := range a.services {
		if b.services[n].ID != svc.ID {
			c.Logf("ID mismatch, a.services[%v]: %v b.services[%v]: %v", n, svc.ID, n, b.services[n].ID)
			sameVals = false
		}
		if b.services[n].DesiredState != svc.DesiredState {
			c.Logf("DesiredState mismatch, a.services[%v]: %v b.services[%v]: %v", n, svc.DesiredState, n, b.services[n].DesiredState)
			sameVals = false
		}
	}
	return sameVals
}

func (s *ServiceStateManagerSuite) LogBatch(c *C, b ServiceStateChangeBatch) {
	svcStr := ""
	for _, svc := range b.services {
		svcStr += fmt.Sprintf(`&service.Service{
			ID: %v,
			DesiredState: %v,
			EmergencyShutdownLevel: %v,
			StartLevel: %v,
		},
		`, svc.ID, svc.DesiredState, svc.EmergencyShutdownLevel, svc.StartLevel)
	}

	c.Logf(`ServiceStateChangeBatch{
	services: []*service.Service{
		%v
	},
	desiredState: %v,
	emergency: %v,
}`, svcStr, b.desiredState, b.emergency)
}
