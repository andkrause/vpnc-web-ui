export interface Status {
  currentPublicIp: string;
  activeVpnClient: string;
  activeVpnConfig: string;
  message: string;
}

export interface ConnectionStatus {
  isActive: boolean;
}

export interface VPNConfig {
  id: string;
  vpnClientName: string;
  configName: string;
  status: ConnectionStatus;
}

export interface DesiredConnectionStatus {
  desiredConnectionState: 'active' | 'inactive';
}

export interface ApiError {
  code: string;
  message: string;
} 