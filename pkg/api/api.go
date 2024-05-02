package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andkrause/vpnc-web-ui/gen/vpnapi"
	"github.com/andkrause/vpnc-web-ui/pkg/vpnclient"
)

// VpnConnectionApiService is a service that implements the logic for the VpnConnectionApiServicer
// This service should implement the business logic for every endpoint for the VpnConnectionApi API.
// Include any external packages or services that will be required by this service.
type ApiService struct {
	vpnClientAggregator *vpnclient.VpnClientAggregator
}

// NewVpnConnectionApiService creates a default api service
func New(vpnClientAggregator *vpnclient.VpnClientAggregator) *ApiService {
	return &ApiService{
		vpnClientAggregator: vpnClientAggregator,
	}
}

func (s *ApiService) OverallStatus(ctx context.Context) (vpnapi.ImplResponse, error) {
	vpnStatus := s.vpnClientAggregator.Status()

	return vpnapi.Response(http.StatusOK, vpnapi.Status{
		CurrentPublicIp: vpnStatus.CurrentPublicIp,
		ActiveVpnClient: vpnStatus.ActiveVpnClient,
		ActiveVpnConfig: vpnStatus.ActiveVpnConfig,
		Message:         vpnStatus.Message,
	}), nil
}

// ListConnections - Get list of Connections
func (s *ApiService) ListConnections(ctx context.Context) (vpnapi.ImplResponse, error) {
	configurations, err := s.vpnClientAggregator.ConfigurationList()
	if err != nil {
		return vpnapi.Response(http.StatusBadRequest, vpnapi.Error{
			Message: fmt.Sprintf("Error reading available connections: %s", err.Error()),
			Code:    "confListUnavailable",
		}), nil
	}

	vpnStatus := s.vpnClientAggregator.Status()

	response := make([]vpnapi.VpnConfig, len(configurations))

	for i := range configurations {
		response[i] = vpnapi.VpnConfig{
			Id:            configurations[i].VpnClientName + configurations[i].VPNConfigurationName,
			VpnClientName: configurations[i].VpnClientName,
			ConfigName:    configurations[i].VPNConfigurationName,
		}
		if vpnStatus.ActiveVpnClient == configurations[i].VpnClientName &&
			vpnStatus.ActiveVpnConfig == configurations[i].VPNConfigurationName {
			response[i].Status = vpnapi.ConnectionStatus{IsActive: true}
		} else {
			response[i].Status = vpnapi.ConnectionStatus{IsActive: false}
		}
	}

	return vpnapi.Response(http.StatusOK, response), nil
}

// ReadConnectionStatus - Read connection status
func (s *ApiService) ReadConnectionStatus(ctx context.Context, id string, client string) (vpnapi.ImplResponse, error) {

	if !s.vpnClientAggregator.ConfigurationExists(client, id) {
		return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
			Message: fmt.Sprintf("VPN configuration %q does not exist", id),
			Code:    "connectionNotFound",
		}), nil
	}

	status := s.vpnClientAggregator.Status()

	if status.ActiveVpnConfig == id && status.ActiveVpnClient == client {
		return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
			IsActive: true,
		}), nil
	}

	return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
		IsActive: false,
	}), nil
}

// SetConnectionStatus - Set connection status
func (s *ApiService) SetConnectionStatus(ctx context.Context, id string, client string, desiredConnectionStatus vpnapi.DesiredConnectionStatus) (vpnapi.ImplResponse, error) {

	if !s.vpnClientAggregator.ConfigurationExists(client, id) {
		return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
			Message: fmt.Sprintf("VPN configuration %q does not exist", id),
			Code:    "connectionNotFound",
		}), nil
	}

	if desiredConnectionStatus.DesiredConnectionState == "active" {

		if err := s.vpnClientAggregator.Connect(client, id); err != nil {
			return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
				Message: fmt.Sprintf("Error establishing VPN connection: %s", err.Error()),
				Code:    "connectionError",
			}), nil
		}
		return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
			IsActive: true,
		}), nil
	} else {

		if err := s.vpnClientAggregator.Disconnect(); err != nil {
			return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
				Message: fmt.Sprintf("Error disconnecting VPN connection: %s", err.Error()),
				Code:    "disconnectingError",
			}), nil
		}
		return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
			IsActive: false,
		}), nil
	}
}
