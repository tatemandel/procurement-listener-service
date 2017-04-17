// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inmemory

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"procurementlistenerservice/model"
	"reflect"
)

type EntitlementState int

const (
	ACTIVE EntitlementState = iota
)

// EntitlementInfo is the internal state that this service holds about an entitlement.
type EntitlementInfo struct {
	Id        string
	State     EntitlementState
	ServiceId string
	PlanId    string

	AccountId   string
	RequestorId string
	Parameters  map[string]interface{}
}

type InMemoryService struct {
	Metadata     Metadata
	Entitlements map[string]EntitlementInfo
}

var _ model.PartnerBackendService = &InMemoryService{}

// CreateService creates a new InMemoryService and returns.
func CreateService(metadata Metadata) *InMemoryService {
	return &InMemoryService{
		Metadata:     metadata,
		Entitlements: make(map[string]EntitlementInfo),
	}
}

// Reset clears the in-memory state.
func (s *InMemoryService) Reset() {
	s.Entitlements = make(map[string]EntitlementInfo)
}

func (s *InMemoryService) OnEntitlementEvent(e model.EntitlementEvent) (model.EntitlementEventResponse, error) {
	switch e.EventType {
	case model.ENTITLEMENT_CREATED:
		return s.onEntitlementCreated(e)
	}

	return model.EntitlementEventResponse{}, fmt.Errorf("Unrecognized entitlement event: '%+v'", e)
}

func (s *InMemoryService) onEntitlementCreated(e model.EntitlementEvent) (model.EntitlementEventResponse, error) {
	serviceDef, err := s.Metadata.getService(e.ServiceId)
	if err != nil {
		log.Printf("Service not found: '%s'.", e.ServiceId)
		return model.EntitlementEventResponse{
			Status: model.RESPONSESTATUS_INVALIDREQUEST,
		}, nil
	}

	planDef, err := serviceDef.getPlan(e.PlanId)
	if err != nil {
		log.Printf("Plan not found: '%s'.", e.PlanId)
		return model.EntitlementEventResponse{
			Status: model.RESPONSESTATUS_INVALIDREQUEST,
		}, nil
	}

	err = validateParameters(e.Parameters, planDef.InputParameterSchema)
	if err != nil {
		log.Printf("Parameters are not valid: '%+v'", err)
		return model.EntitlementEventResponse{
			Status: model.RESPONSESTATUS_INVALIDREQUEST,
		}, nil
	}

	state := EntitlementInfo{
		Id:          e.EntitlementId,
		State:       ACTIVE,
		ServiceId:   e.ServiceId,
		PlanId:      e.PlanId,
		AccountId:   e.AccountId,
		RequestorId: e.RequestorId,
		Parameters:  e.Parameters,
	}

	existing, exists := s.Entitlements[e.EntitlementId]
	if exists {
		if !reflect.DeepEqual(existing, state) {
			log.Printf("Entitlement already exists: '%s'.", e.EntitlementId)
			return model.EntitlementEventResponse{
				Status: model.RESPONSESTATUS_INVALIDREQUEST,
			}, nil
		}
	} else {
		s.Entitlements[e.EntitlementId] = state
	}

	log.Printf("Entitlement created: '%+v'\n", state)

	return model.EntitlementEventResponse{
		Status:  model.RESPONSESTATUS_ACCEPTED,
		EventId: e.EventId,
	}, nil
}

func validateParameters(parameters map[string]interface{}, schema map[string]interface{}) error {
	if len(schema) == 0 {
		// No schema was defined
		if len(parameters) != 0 {
			return errors.New("No parameters were expected.")
		}
		return nil
	}

	parametersLoader := gojsonschema.NewGoLoader(parameters)
	schemaLoader := gojsonschema.NewGoLoader(schema)

	result, err := gojsonschema.Validate(schemaLoader, parametersLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return fmt.Errorf("The document is not valid: {%v}.", result.Errors())
	}

	return nil
}
