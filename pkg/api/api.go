package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andkrause/vpnc-web-ui/gen/vpnapi"
	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"
)

// VpnConnectionApiService is a service that implements the logic for the VpnConnectionApiServicer
// This service should implement the business logic for every endpoint for the VpnConnectionApi API.
// Include any external packages or services that will be required by this service.
type ApiService struct {
	vpnc *vpnc.VPNC
}

// NewVpnConnectionApiService creates a default api service
func New(vpnc *vpnc.VPNC) *ApiService {
	return &ApiService{
		vpnc: vpnc,
	}
}

func (s *ApiService) OverallStatus(ctx context.Context) (vpnapi.ImplResponse, error) {
	vpnStatus := s.vpnc.Status()

	return vpnapi.Response(http.StatusOK, vpnapi.Status{
		CurrentPublicIp: vpnStatus.CurrentPublicIp,
		ActiveVpnConfig: vpnStatus.ActiveVpnConfig,
		Message:         vpnStatus.Message,
	}), nil
}

// ListConnections - Get list of Connections
func (s *ApiService) ListConnections(ctx context.Context) (vpnapi.ImplResponse, error) {
	configurations, err := s.vpnc.ConfigurationList()
	if err != nil {
		return vpnapi.Response(http.StatusBadRequest, vpnapi.Error{
			Message: fmt.Sprintf("Error reading available connections: %s", err.Error()),
			Code:    "confListUnavailable",
		}), nil
	}

	vpnStatus := s.vpnc.Status()

	response := make([]vpnapi.VpnConfig, len(configurations))

	for i, conf := range configurations {

		response[i] = vpnapi.VpnConfig{
			Id:   conf,
			Name: conf,
		}

		if conf == vpnStatus.ActiveVpnConfig {
			response[i].Status = vpnapi.ConnectionStatus{IsActive: true}
		} else {
			response[i].Status = vpnapi.ConnectionStatus{IsActive: false}
		}
	}

	return vpnapi.Response(http.StatusOK, response), nil
}

// ReadConnectionStatus - Read connection status
func (s *ApiService) ReadConnectionStatus(ctx context.Context, id string) (vpnapi.ImplResponse, error) {

	if !s.vpnc.ConfigurationExists(id) {
		return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
			Message: fmt.Sprintf("VPN configuration %q does not exist", id),
			Code:    "connectionNotFound",
		}), nil
	}

	if s.vpnc.Status().ActiveVpnConfig == id {
		return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
			IsActive: true,
		}), nil
	}

	return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
		IsActive: false,
	}), nil
}

// SetConnectionStatus - Set connection status
func (s *ApiService) SetConnectionStatus(ctx context.Context, id string, desiredConnectionStatus vpnapi.DesiredConnectionStatus) (vpnapi.ImplResponse, error) {

	if !s.vpnc.ConfigurationExists(id) {
		return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
			Message: fmt.Sprintf("VPN configuration %q does not exist", id),
			Code:    "connectionNotFound",
		}), nil
	}

	if desiredConnectionStatus.DesiredConnectionState == "active" {
		if err := s.vpnc.Connect(id); err != nil {
			return vpnapi.Response(http.StatusNotFound, vpnapi.Error{
				Message: fmt.Sprintf("Error establishing VPN connection: %s", err.Error()),
				Code:    "connectionError",
			}), nil
		}
		return vpnapi.Response(http.StatusOK, vpnapi.ConnectionStatus{
			IsActive: true,
		}), nil
	} else {
		if err := s.vpnc.Disconnect(); err != nil {
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
