/*
 * VPN Gateway API
 *
 * A REST compliant API to manage a VPN Gateway instance.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package vpnapi

type DesiredConnectionStatus struct {

	// Indicates whether a specific should be active or not
	DesiredConnectionState string `json:"desiredConnectionState,omitempty"`
}

// AssertDesiredConnectionStatusRequired checks if the required fields are not zero-ed
func AssertDesiredConnectionStatusRequired(obj DesiredConnectionStatus) error {
	return nil
}

// AssertRecurseDesiredConnectionStatusRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DesiredConnectionStatus (e.g. [][]DesiredConnectionStatus), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDesiredConnectionStatusRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDesiredConnectionStatus, ok := obj.(DesiredConnectionStatus)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDesiredConnectionStatusRequired(aDesiredConnectionStatus)
	})
}