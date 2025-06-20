// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * VPN Gateway API
 *
 * A REST compliant API to manage a VPN Gateway instance.
 *
 * API version: 1.0
 */

package vpnapi




type VpnConfig struct {

	// ID of a VPN Configuration that the gateway can connect to
	Id string `json:"id"`

	// Name of the VPN Client
	VpnClientName string `json:"vpnClientName"`

	// Human readable name of a VPN Configuration that the gateway can connect to
	ConfigName string `json:"configName"`

	Status ConnectionStatus `json:"status"`
}

// AssertVpnConfigRequired checks if the required fields are not zero-ed
func AssertVpnConfigRequired(obj VpnConfig) error {
	elements := map[string]interface{}{
		"id": obj.Id,
		"vpnClientName": obj.VpnClientName,
		"configName": obj.ConfigName,
		"status": obj.Status,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertConnectionStatusRequired(obj.Status); err != nil {
		return err
	}
	return nil
}

// AssertVpnConfigConstraints checks if the values respects the defined constraints
func AssertVpnConfigConstraints(obj VpnConfig) error {
	if err := AssertConnectionStatusConstraints(obj.Status); err != nil {
		return err
	}
	return nil
}
