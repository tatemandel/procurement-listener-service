package inmemory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Metadata is the top-level container of metadata.
type Metadata struct {
	Services []ServiceDefinition `json:"services"`
}

// ServiceDefinition is the metadata about a particular service that this procurement backend handles.
type ServiceDefinition struct {
	// ServiceId is the id of the service.
	ServiceId string `json:"serviceId"`

	Plans []PlanDefinition `json:"plans"`
}

// PlanDefinition is the metadata about a particular plan that this procurement backend handles.
type PlanDefinition struct {
	// PlanId is the id of the plan.
	PlanId               string                 `json:"planId"`
	InputParameterSchema map[string]interface{} `json:"inputParameterSchema"`
}

func (m *Metadata) getService(id string) (ServiceDefinition, error) {
	for _, def := range m.Services {
		if def.ServiceId == id {
			return def, nil
		}
	}
	return ServiceDefinition{}, fmt.Errorf("ServiceDefinition not found: id='%s'.", id)
}

func (s *ServiceDefinition) getPlan(id string) (PlanDefinition, error) {
	for _, def := range s.Plans {
		if def.PlanId == id {
			return def, nil
		}
	}
	return PlanDefinition{}, fmt.Errorf("PlanDefinition not found: id='%s'.", id)
}

// ReadMetadataFile opens the file with the given path, reads contents as JSON, and returns the parsed Metadata struct.
func ReadMetadataFile(path string) (Metadata, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Metadata{}, fmt.Errorf("Unable to read metadata file: '%s'.\n", path)
	}

	var metadata Metadata
	err = json.Unmarshal(contents, &metadata)
	if err != nil {
		return Metadata{}, fmt.Errorf("Unable to parse metadata file: '%v'.\n", err)
	}

	return metadata, nil
}
