// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * VPN Gateway API
 *
 * A REST compliant API to manage a VPN Gateway instance.
 *
 * API version: 1.0
 */

package vpnapi

type ConnectionStatus struct {

	// Indicates whether a specific connection is currently active or not
	IsActive bool `json:"isActive"`
}

// AssertConnectionStatusRequired checks if the required fields are not zero-ed
func AssertConnectionStatusRequired(obj ConnectionStatus) error {
	elements := map[string]interface{}{
		"isActive": obj.IsActive,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertConnectionStatusConstraints checks if the values respects the defined constraints
func AssertConnectionStatusConstraints(obj ConnectionStatus) error {
	return nil
}
