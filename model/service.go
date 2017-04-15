package model

import "fmt"

// EntitlementEventType is the underlying type for entitlement event related enum values.
type EntitlementEventType string

const (
	// ENTITLEMENT_CREATED indicates that an entitlement has been created on the source system.
	ENTITLEMENT_CREATED EntitlementEventType = "ENTITLEMENT_CREATED"

	// ENTITLEMENT_DELETED indicates that an entitlement has been deleted on the source system.
	ENTITLEMENT_DELETED EntitlementEventType = "ENTITLEMENT_DELETED"

	// ENTITLEMENT_UPDATED indicates that an entitlement has been updated on the source system.
	ENTITLEMENT_UPDATED EntitlementEventType = "ENTITLEMENT_UPDATED"

	// ENTITLEMENT_CANCELLED indicates that an entitlement has been cancelled by the owner on the source system.
	ENTITLEMENT_CANCELLED EntitlementEventType = "ENTITLEMENT_CANCELLED"

	// ENTITLEMENT_REACTIVATED indicates that an entitlement has been reactivated by the owner on the source system.
	ENTITLEMENT_REACTIVATED EntitlementEventType = "ENTITLEMENT_REACTIVATED"
)

// EntitlementEvent is received when an entitlement related notification event is sent from the source system
// to the service provider backend.
//
// An Entitlement event is uniquely identified by EventId. Due to the unreliable nature of networks, the same event
// (with the same EventId) can be dispatched multiple times.
//
type EntitlementEvent struct {

	// EventId uniquely identifies this event.
	EventId string `json:"eventId"`

	// EventType is the type of this event. Based on the EventType, certain fields might be missing.
	EventType EntitlementEventType `json:"eventType"`

	// ServiceId identifies the service that the entitlement is for.
	// This is set when EventId is ENTITLEMENT_CREATED.
	ServiceId string `json:"serviceId"`

	// PlanId identifies the particular plan that was chosen by the user during the entitlement creation.
	// This is set when EventId is ENTITLEMENT_CREATED.
	PlanId string `json:"planId"`

	// Entitlement is the full URL to the entitlement resource that this event is based on.
	EntitlementId string `json:"entitlementId"`

	// AccountId is the account id of entitlement's owner in the marketplace.
	AccountId string `json:"accountId"`

	// RequestorId is the account id of an intermediary, who triggered the entitlement event on behalf of the owner,
	// if any.
	RequestorId string `json:"requestorId"`

	// Parameters are the custom parameters that was supplied by the user, as part of this event.
	Parameters map[string]interface{} `json:"parameters"`
}

// ResponseStatus is the underlying type for the response status relate enum values
type ResponseStatus int

const (
	// RESPONSESTATUS_INVALIDREQUEST indicates that the request was invalid.
	RESPONSESTATUS_INVALIDREQUEST ResponseStatus = iota

	// RESPONSESTATUS_ACCEPTED indicates that the event was accepted.
	RESPONSESTATUS_ACCEPTED ResponseStatus = iota

	// RESPONSESTATUS_REJECTED indicates that the event was rejected.
	RESPONSESTATUS_REJECTED ResponseStatus = iota

	// RESPONSESTATUS_ASYNC indicates that the event will be accepted/rejected asynchronously.
	RESPONSESTATUS_ASYNC ResponseStatus = iota
)

// EntitlementEventResponse represents a response to an entitlement event notification that was received from the source system.
type EntitlementEventResponse struct {

	// The response status for the
	Status ResponseStatus `json:-`

	// eventId is the id of the event that this response is being returned for.
	EventId string `json:"eventId"`

	// entitlementDashboardUrl is a templatized SSO dashboard url that the entitlement owner can use to manage
	// the entitlement on the service provider side.
	EntitlementDashboardUrl string `json:"entitlementDashboardUrl"`

	// labels are optional custom parameters that the backend would like to attach to the entitlement.
	Labels struct{} `json:"labels"`
}

// PartnerBackendService is the service interface that needs to be implemented by the backends to listen and react to incoming procurement events.
type PartnerBackendService interface {
	// TODO: Add support for async handling of the events.
	// OnEntitlementvents gets invoked when a new entitlement event is received.
	OnEntitlementEvent(e EntitlementEvent) (EntitlementEventResponse, error)
}

func ValidateEntitlementEvent(e EntitlementEvent) error {
	if e.EventId == "" {
		return fmt.Errorf("Field 'eventId' does not have a valid value: '%v'.", e.EventId)
	}

	if e.EntitlementId == "" {
		return fmt.Errorf("Field 'entitlementId' does not have a valid value: '%v'.", e.EntitlementId)
	}

	switch e.EventType {
	case ENTITLEMENT_CREATED:
		if e.ServiceId == "" {
			return fmt.Errorf("Field 'serviceId' does not have a valid value: '%v'.", e.ServiceId)
		}
		if e.PlanId == "" {
			return fmt.Errorf("Field 'planId' does not have a valid value: '%v'.", e.PlanId)
		}
		break
	case ENTITLEMENT_DELETED:
	case ENTITLEMENT_UPDATED:
	case ENTITLEMENT_CANCELLED:
	case ENTITLEMENT_REACTIVATED:
		break
	default:
		return fmt.Errorf("Field 'eventType' doesn't have a valid value: '%v'.", e.EventType)
	}

	return nil
}
