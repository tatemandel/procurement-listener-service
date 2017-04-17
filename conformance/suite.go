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

//
// Package conformance contains a basic conformance test suite for the procurement listener service, along with the
// infrastructure to drive the conformance tests.
//
package conformance

import (
	"encoding/json"
	"log"
	"procurementlistenerservice/inmemory"
)

// This files contains behavior/conformance tests for the Procurement Partner Backend, in declarative form.

// Metadata for the test setup.
var metadata inmemory.Metadata = inmemory.Metadata{
	Services: []inmemory.ServiceDefinition{
		{
			// A Simple service with a single plan, with no inputs expected.
			ServiceId: "Simple",
			Plans: []inmemory.PlanDefinition{
				{
					PlanId: "SimplePlan1",
				},
			},
		},

		{
			// A more complicated service that expects parameters as part of creation.
			ServiceId: "Parameterized",
			Plans: []inmemory.PlanDefinition{
				{
					PlanId: "ParameterizedPlan1",
					InputParameterSchema: createInputParameterSchema(`
					{
					    "title": "SimpleParameterized Input Schema",
					    "type": "object",
					    "properties": {
					      "parameter1": {
						"type": "string"
					      },
					      "parameter2": {
						"type": "integer",
						"minimum": 0
					      }
					    },
					    "required": ["parameter2"]
					}
					`),
				},
			},
		},
	},
}

func createInputParameterSchema(contents string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(contents), &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// Tests include all conformance tests to be run against a Procurement Listener Service.
var Tests []ConformanceTest = []ConformanceTest{
	{
		Name:     "simpleSuccess",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Simple",
						PlanId:    "SimplePlan1",
						State:     inmemory.ACTIVE,
					},
				},
			},
		},
	},

	{
		Name:     "emptyRequest",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
				}`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "unknownService",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "WorldDominationService",
					"planId": "trial"
				}`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleMissingEventType",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleMissingEventId",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleMissingServiceId",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleMissingPlanId",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleRedundantParameters",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1",
					"parameters": {
						"foo": "bar"
					}
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleUnknownPlan",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "WorldDomination"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "simpleDoubleCreate",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Simple",
						PlanId:    "SimplePlan1",
						State:     inmemory.ACTIVE,
					},
				},
			},
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Simple",
						PlanId:    "SimplePlan1",
						State:     inmemory.ACTIVE,
					},
				},
			},
		},
	},

	{
		Name:     "simpleDuplicateEvent",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Simple",
						PlanId:    "SimplePlan1",
						State:     inmemory.ACTIVE,
					},
				},
			},
		},
	},

	{
		Name:     "simpleDuplicateEvent",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Simple",
					"planId": "SimplePlan1"
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Simple",
						PlanId:    "SimplePlan1",
						State:     inmemory.ACTIVE,
					},
				},
			},
		},
	},

	{
		Name:     "parameterizedEmptyParameters",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Parameterized",
					"planId": "ParameterizedPlan1"
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "parameterizedInvalidParameters",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Parameterized",
					"planId": "ParameterizedPlan1",
					"parameters": {
						"foo": "bar",
					}
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "parameterizedRequiredParameterMissing",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Parameterized",
					"planId": "ParameterizedPlan1",
					"parameters": {
						"parameter1": "exists",
					}
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "parameterizedRequiredParameterMissing",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Parameterized",
					"planId": "ParameterizedPlan1",
					"parameters": {
						"parameter1": "exists",
					}
				}
				`,
				ExpectedCode: 400,
			},
		},
	},

	{
		Name:     "parameterizedSuccess",
		Metadata: metadata,
		Actions: []Action{
			PostEntitlementEvent{
				Request: `
				{
					"eventId": "1",
					"eventType": "ENTITLEMENT_CREATED",
					"entitlementId": "E1",
					"serviceId": "Parameterized",
					"planId": "ParameterizedPlan1",
					"parameters": {
						"parameter2": 42
					}
				}
				`,
				ExpectedCode: 200,
			},
			ExpectEntitlements{
				Entitlements: []inmemory.EntitlementInfo{
					{
						Id:        "E1",
						ServiceId: "Parameterized",
						PlanId:    "ParameterizedPlan1",
						State:     inmemory.ACTIVE,
						Parameters: map[string]interface{}{
							"parameter2": 42.,
						},
					},
				},
			},
		},
	},
}
