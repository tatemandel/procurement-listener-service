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

var tests []ConformanceTest = []ConformanceTest{
	emptyRequest(),
	unknownService(),

	simpleMissingEventType(),
	simpleMissingEventId(),
	simpleMissingServiceId(),
	simpleMissingPlanId(),
	simpleRedundantParameters(),
	simpleSuccess(),
	simpleUnknownPlan(),
	simpleDoubleCreate(),
	simpleDuplicateEvent(),

	parameterizedEmptyParameters(),
	parameterizedInvalidParameters(),
	parameterizedRequiredParameterMissing(),
	parameterizedSuccess(),
}

func emptyRequest() ConformanceTest {
	b := BeginTestCase("emptyRequest", metadata)

	b.PostEntitlementEvent(`
	{
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func unknownService() ConformanceTest {
	b := BeginTestCase("unknownService", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "WorldDominationService",
		"planId": "trial"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleMissingEventType() ConformanceTest {
	b := BeginTestCase("simpleMissingEventType", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleMissingEventId() ConformanceTest {
	b := BeginTestCase("simpleMissingEventId", metadata)

	b.PostEntitlementEvent(`
	{
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleMissingServiceId() ConformanceTest {
	b := BeginTestCase("simpleMissingServiceId", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleMissingPlanId() ConformanceTest {
	b := BeginTestCase("simpleMissingPlanId", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleRedundantParameters() ConformanceTest {
	b := BeginTestCase("simpleRedundantParameters", metadata)

	b.PostEntitlementEvent(`
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
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleUnknownPlan() ConformanceTest {
	b := BeginTestCase("simpleUnknownPlan", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "WorldDomination"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func simpleDoubleCreate() ConformanceTest {
	b := BeginTestCase("simpleDoubleCreate", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(200)

	b.ExpectEntitlementCount(1)

	b.PostEntitlementEvent(`
	{
		"eventId": "2",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(200)

	b.ExpectEntitlementCount(1)
	b.ExpectEntitlement(inmemory.EntitlementInfo{
		Id:        "E1",
		ServiceId: "Simple",
		PlanId:    "SimplePlan1",
		State:     inmemory.ACTIVE,
	})

	return b.Build()
}

func simpleDuplicateEvent() ConformanceTest {
	b := BeginTestCase("simpleDuplicateEvent", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(200)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(200)

	return b.Build()
}

func simpleSuccess() ConformanceTest {
	b := BeginTestCase("simpleSuccess", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Simple",
		"planId": "SimplePlan1"
	}
	`).ExpectResponseCode(200)

	b.ExpectEntitlementCount(1)
	b.ExpectEntitlement(inmemory.EntitlementInfo{
		Id:        "E1",
		ServiceId: "Simple",
		PlanId:    "SimplePlan1",
		State:     inmemory.ACTIVE,
	})

	return b.Build()
}

func parameterizedEmptyParameters() ConformanceTest {
	b := BeginTestCase("parameterizedEmptyParameters", metadata)

	b.PostEntitlementEvent(`
	{
		"eventId": "1",
		"eventType": "ENTITLEMENT_CREATED",
		"entitlementId": "E1",
		"serviceId": "Parameterized",
		"planId": "ParameterizedPlan1"
	}
	`).ExpectResponseCode(400)

	return b.Build()
}

func parameterizedInvalidParameters() ConformanceTest {
	b := BeginTestCase("parameterizedInvalidParameters", metadata)

	b.PostEntitlementEvent(`
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
	`).ExpectResponseCode(400)

	return b.Build()
}

func parameterizedRequiredParameterMissing() ConformanceTest {
	b := BeginTestCase("parameterizedRequiredParameterMissing", metadata)

	b.PostEntitlementEvent(`
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
	`).ExpectResponseCode(400)

	return b.Build()
}

func parameterizedSuccess() ConformanceTest {
	b := BeginTestCase("parameterizedSuccess", metadata)

	b.PostEntitlementEvent(`
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
	`).ExpectResponseCode(200)

	b.ExpectEntitlementCount(1)
	b.ExpectEntitlement(inmemory.EntitlementInfo{
		Id:        "E1",
		ServiceId: "Parameterized",
		PlanId:    "ParameterizedPlan1",
		State:     inmemory.ACTIVE,
		Parameters: map[string]interface{}{
			"parameter2": 42.,
		},
	})

	return b.Build()
}

func createInputParameterSchema(contents string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(contents), &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
