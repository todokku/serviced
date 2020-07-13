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
	"fmt"

	"github.com/control-center/serviced/dao"
	"github.com/control-center/serviced/domain"
	// "github.com/control-center/serviced/domain/logfilter"
	"github.com/control-center/serviced/domain/service"
	definition "github.com/control-center/serviced/domain/servicedefinition"
	"github.com/control-center/serviced/utils"
)

func (mc *MigrationContext) validate(req dao.ServiceMigrationRequest) error {
	var svcAll []service.Service

	// Validate service updates
	for _, svc := range req.Modified {
		if _, err := mc.validateUpdate(svc); err != nil {
			return err
		}
		svcAll = append(svcAll, *svc)
	}

	// Validate service adds
	for _, svc := range req.Added {
		if err := mc.validateAdd(svc); err != nil {
			return err
		} else if svc.ID, err = utils.NewUUID36(); err != nil {
			return err
		}
		svcAll = append(svcAll, *svc)
	}

	// Validate service deployments
	for _, sdreq := range req.Deploy {
		err := mc.validateDeployment(sdreq.ParentID, &sdreq.Service)
		if err != nil {
			return err
		}
		// svcAll = append(svcAll, svcs...)
	}

	// Validate service migration
	if err := mc.validateMigration(svcAll, req.ServiceID); err != nil {
		return err
	}

	return nil
}

func (mc *MigrationContext) validateUpdate(svc *service.Service) (*service.Service, error) {

	// Verify that the service exists.
	cursvc, err := mc.services.Get(svc.ID)
	if err != nil {
		return nil, err
	}

	// verify no collision with the service name
	if svc.Name != cursvc.Name {
		if err := mc.validateName(svc); err != nil {
			return nil, err
		}
	}

	// disallow enabling ports and vhosts that are already enabled by a different
	// service and application.
	// TODO: what if they are on the same service?
	// for _, ep := range svc.Endpoints {
	// 	for _, vhost := range ep.VHostList {
	// 		if vhost.Enabled {
	// 			serviceID, application, err := f.zzk.GetVHost(vhost.Name)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			if (serviceID != "" && serviceID != svc.ID) ||
	// 				(application != "" && application != ep.Application) {
	// 				return nil, fmt.Errorf("vhost %s is already in use", vhost.Name)
	// 			}
	// 		}
	// 	}
	// 	for _, port := range ep.PortList {
	// 		if port.Enabled {
	// 			serviceID, application, err := f.zzk.GetPublicPort(port.PortAddr)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			if (serviceID != "" && serviceID != svc.ID) ||
	// 				(application != "" && application != ep.Application) {
	// 				return nil, fmt.Errorf("port %s is already in use", port.PortAddr)
	// 			}
	// 		}
	// 	}
	// }

	// set read-only fields
	svc.CreatedAt = cursvc.CreatedAt
	svc.DeploymentID = cursvc.DeploymentID

	// remove any BuiltIn enabled monitoring configs
	metricConfigs := []domain.MetricConfig{}
	for _, mcfg := range svc.MonitoringProfile.MetricConfigs {
		if mcfg.ID == "metrics" {
			continue
		}
		metrics := []domain.Metric{}
		for _, m := range mcfg.Metrics {
			if !m.BuiltIn {
				metrics = append(metrics, m)
			}
		}
		mcfg.Metrics = metrics
		metricConfigs = append(metricConfigs, mcfg)
	}
	svc.MonitoringProfile.MetricConfigs = metricConfigs

	graphs := []domain.GraphConfig{}
	for _, g := range svc.MonitoringProfile.GraphConfigs {
		if !g.BuiltIn {
			graphs = append(graphs, g)
		}
	}
	svc.MonitoringProfile.GraphConfigs = graphs

	if err := validateServiceOptions(svc); err != nil {
		return nil, err
	}

	return cursvc, nil
}

// Validates that the service doesn't have invalid options specified.
func validateServiceOptions(svc *service.Service) error {
	// ChangeOption RestartAllOnInstanceChanged and HostPolicy RequireSeparate are invalid together.
	var changeOptions = definition.ChangeOptions(svc.ChangeOptions)
	if svc.HostPolicy == definition.RequireSeparate &&
		changeOptions.Contains(definition.RestartAllOnInstanceChanged) {
		return fmt.Errorf(
			"HostPolicy RequireSeparate cannot be used with ChangeOption RestartAllOnInstanceChanged",
		)
	}
	return nil
}

func (mc *MigrationContext) validateName(svc *service.Service) error {
	if svc.ParentServiceID != "" {
		parentSvc, err := mc.services.Get(svc.ParentServiceID)
		if err != nil {
			return err
		}
		svc.DeploymentID = parentSvc.DeploymentID
	}
	cursvc, err := mc.services.FindChild(svc.ParentServiceID, svc.Name)
	if err != nil {
		return err
	}
	if cursvc != nil {
		path, err := mc.services.GetServicePath(svc.Name)
		if err != nil {
			path = fmt.Sprintf("%v", err)
		}
		return fmt.Errorf("Service already exists on path.  path=%s", path)
	}
	return nil
}

func (mc *MigrationContext) validateAdd(svc *service.Service) error {
	// Verify that the service does not already exist
	existing, err := mc.services.Get(svc.ID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("Service exists: name=%s id=%s", svc.Name, svc.ID)
	}

	// verify no collision with the service name
	if err := mc.validateName(svc); err != nil {
		return err
	}

	// disable ports and vhosts that are already in use by another application
	// for i, ep := range svc.Endpoints {
	// 	for j, vhost := range ep.VHostList {
	// 		if vhost.Enabled {
	// 			serviceID, application, err := f.zzk.GetVHost(vhost.Name)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if serviceID != "" || application != "" {
	// 				svc.Endpoints[i].VHostList[j].Enabled = false
	// 			}
	// 		}
	// 	}

	// 	for j, port := range ep.PortList {
	// 		if port.Enabled {
	// 			serviceID, application, err := f.zzk.GetPublicPort(port.PortAddr)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if serviceID != "" || application != "" {
	// 				svc.Endpoints[i].PortList[j].Enabled = false
	// 			}
	// 		}
	// 	}
	// }

	if err := validateServiceOptions(svc); err != nil {
		return err
	}

	// remove any BuiltIn enabled monitoring configs
	metricConfigs := []domain.MetricConfig{}
	for _, mc := range svc.MonitoringProfile.MetricConfigs {
		if mc.ID == "metrics" {
			continue
		}

		metrics := []domain.Metric{}
		for _, m := range mc.Metrics {
			if !m.BuiltIn {
				metrics = append(metrics, m)
			}
		}
		mc.Metrics = metrics
		metricConfigs = append(metricConfigs, mc)
	}
	svc.MonitoringProfile.MetricConfigs = metricConfigs

	graphs := []domain.GraphConfig{}
	for _, g := range svc.MonitoringProfile.GraphConfigs {
		if !g.BuiltIn {
			graphs = append(graphs, g)
		}
	}
	svc.MonitoringProfile.GraphConfigs = graphs

	// set service defaults
	svc.DesiredState = int(service.SVCStop)         // new services must always be stopped
	svc.CurrentState = string(service.SVCCSStopped) // new services are always stopped
	// manage service configurations
	if svc.OriginalConfigs == nil || len(svc.OriginalConfigs) == 0 {
		if svc.ConfigFiles != nil {
			svc.OriginalConfigs = svc.ConfigFiles
		} else {
			svc.OriginalConfigs = make(map[string]definition.ConfigFile)
		}
	}
	return nil
}

// validateDeployment returns the services that will be deployed
func (mc *MigrationContext) validateDeployment(
	parentID string, sd *definition.ServiceDefinition,
) error {
	var err error
	var parent *service.Service

	parent, err = mc.services.Get(parentID)
	if err != nil {
		return err
	}

	return deploy(
		mc.services,
		mc.services.GetTenantID(),
		parent.PoolID,
		parent.DeploymentID,
		parentID,
		false,
		*sd,
	)
}

// validateMigration makes sure there are no collisions with the added/modified services.
func (mc *MigrationContext) validateMigration(svcs []service.Service, tenantID string) error {
	svcParentMapNameMap := make(map[string]map[string]struct{})
	endpointMap := make(map[string]string)
	for _, svc := range svcs {
		// check for name uniqueness within the set of new/modified/deployed services
		if svcNameMap, ok := svcParentMapNameMap[svc.ParentServiceID]; ok {
			if _, ok := svcNameMap[svc.Name]; ok {
				return fmt.Errorf(
					"Collision for service name %s and parent %s", svc.Name, svc.ParentServiceID,
				)
			}
			svcParentMapNameMap[svc.ParentServiceID][svc.Name] = struct{}{}
		} else {
			svcParentMapNameMap[svc.ParentServiceID] = make(map[string]struct{})
		}

		// check for endpoint name uniqueness within the set of new/modified/deployed services
		for _, ep := range svc.Endpoints {
			if ep.Purpose == "export" {
				if _, ok := endpointMap[ep.Application]; ok {
					return fmt.Errorf(
						"Endpoint %s in migrated service %s is a duplicate of an endpoint in one "+
							"of the other migrated services",
						ep.Application, svc.Name,
					)
				}
				endpointMap[ep.Application] = svc.ID
			}
		}
	}

	// Check whether migrated services' endpoints conflict with existing services.
	var iter = mc.services.Iterator()
	for iter.Next() {
		svc := iter.Item()
		for _, ep := range svc.Endpoints {
			if ep.Purpose != "export" {
				continue
			}
			newsvcID, ok := endpointMap[ep.Application]
			if !ok {
				continue
			}
			if newsvcID != svc.ID {
				return fmt.Errorf(
					"Endpoint %s in migrated service %s is a duplicate of an endpoint in one "+
						"of the other migrated services",
					ep.Application, svc.Name,
				)
			}
		}
	}
	return nil
}
