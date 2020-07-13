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
	// log "github.com/Sirupsen/logrus"
	"fmt"
	"time"

	"github.com/control-center/serviced/dao"
	// "github.com/control-center/serviced/domain"
	"github.com/control-center/serviced/domain/logfilter"
	"github.com/control-center/serviced/domain/service"
	definition "github.com/control-center/serviced/domain/servicedefinition"
	// "github.com/control-center/serviced/utils"
)

// MigrationContext is an entity for migrating services
type MigrationContext struct {
	services *ServiceList
	filters  map[string]logfilter.LogFilter
}

// NewMigrationContext returns a new MigrationContext entity.
func NewMigrationContext(services *ServiceList, filters map[string]logfilter.LogFilter) *MigrationContext {
	return &MigrationContext{
		services: services,
		filters:  filters,
	}
}

// Migrate performs a batch migration on a group of services.
func (mc *MigrationContext) Migrate(req dao.ServiceMigrationRequest) (ServiceList, error) {

	if err := mc.validate(req); err != nil {
		return nil, err
	}

	// Do migration
	// for _, filter := range req.LogFilters {
	// 	var action string
	// 	existingFilter, err := f.logFilterStore.Get(ctx, filter.Name, filter.Version)
	// 	if err == nil {
	// 		existingFilter.Filter = filter.Filter
	// 		err = f.logFilterStore.Put(ctx, existingFilter)
	// 		action = "update"
	// 	} else if err != nil && datastore.IsErrNoSuchEntity(err) {
	// 		err = f.logFilterStore.Put(ctx, &filter)
	// 		action = "add"
	// 	}
	// 	if err != nil {
	// 		logger.WithError(err).WithFields(log.Fields{
	// 			"action":     action,
	// 			"filtername": filter.Name,
	// 		}).Error("Failed to save log filter")
	// 		return err
	// 	}
	// 	logger.WithFields(log.Fields{
	// 		"action":        action,
	// 		"filtername":    filter.Name,
	// 		"filterversion": filter.Version,
	// 	}).Debug("Service migration saved LogFilter")
	// }

	for _, svc := range req.Modified {
		if err := mc.Update(svc); err != nil {
			return nil, err
		}
	}
	for _, svc := range req.Added {
		if err := mc.Add(svc); err != nil {
			return nil, err
		}
	}
	for _, sdreq := range req.Deploy {
		// if err := mc.Deploy(
		// 	"default", "test", sdreq.ParentID, false, sdreq.Service,
		// ); err != nil {
		if err := mc.Deploy(sdreq); err != nil {
			return nil, err
		}
	}
	return *mc.services, nil
}

// Update migrates an existing service; return error if the service does not exist
func (mc *MigrationContext) Update(svc *service.Service) error {
	var tenantID string

	if svc.ParentServiceID == "" {
		tenantID = svc.ID
	} else {
		tenantID = mc.services.GetTenantID()
		if tenantID == "" {
			return fmt.Errorf("No tenant ID found")
		}
	}

	// _, err := mc.validateUpdate(svc)
	// if err != nil {
	// 	return err
	// }

	// set service configurations
	if svc.OriginalConfigs == nil || len(svc.OriginalConfigs) == 0 {
		if svc.ConfigFiles != nil {
			svc.OriginalConfigs = svc.ConfigFiles
		} else {
			svc.OriginalConfigs = make(map[string]definition.ConfigFile)
		}
	}

	svc.UpdatedAt = time.Now()

	mc.services.Put(*svc)
	return nil
}

// Add adds a service; return error if service already exists
func (mc *MigrationContext) Add(svc *service.Service) error {
	// service add validation
	if err := mc.validateAdd(svc); err != nil {
		return err
	}

	svc.UpdatedAt = time.Now()
	svc.CreatedAt = svc.UpdatedAt
	svc.DesiredState = int(service.SVCStop)

	mc.services.Append(*svc)

	return nil
}

// Deploy converts a service definition to a service and deploys it under a specific service.
// If the overwrite option is enabled, existing services with the same name will be overwritten,
// otherwise services may only be added.
// func (mc *MigrationContext) Deploy(
// 	poolID, deploymentID, parentID string, overwrite bool, svcDef definition.ServiceDefinition,
// ) error {
func (mc *MigrationContext) Deploy(req *dao.ServiceDeploymentRequest) error {
	var parent *service.Service
	var err error

	// Get the tenant ID (this is also the ID of the root service)
	tenantID := mc.services.GetTenantID()
	if tenantID == "" {
		return fmt.Errorf("No tenant ID found")
	}

	if req.ParentID == "" {
		return fmt.Errorf("No parent service ID specified")
	}

	// Get the parent service
	parent, err = mc.services.Get(req.ParentID)
	if err != nil {
		return err
	}

	// Do some pool validation
	var poolID = parent.PoolID
	if req.PoolID != "" {
		poolID = req.PoolID
	}
	// if poolID == "" {
	// 	poolID = parent.PoolID
	// }

	// if deploymentID == "" {
	// 	deploymentID = parent.DeploymentID
	// }

	// svc, err := mc.services.FindChild(req.ParentID, req.Service.Name)
	// if err != nil {
	// 	return err
	// }
	// if svc != nil {
	// }

	return deploy(
		mc.services,
		tenantID,
		poolID,
		parent.DeploymentID,
		parent.ID,
		true,
		req.Service,
	)
}
