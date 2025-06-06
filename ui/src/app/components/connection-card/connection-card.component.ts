import { Component, Input, Output, EventEmitter } from '@angular/core';
import { VPNConfig } from '../../models/api.models';

@Component({
  selector: 'app-connection-card',
  templateUrl: './connection-card.component.html',
  styleUrls: ['./connection-card.component.scss']
})
export class ConnectionCardComponent {
  @Input() connection!: VPNConfig;
  @Input() isLoading: boolean = false;
  @Output() connect = new EventEmitter<VPNConfig>();
  @Output() disconnect = new EventEmitter<VPNConfig>();

  get isActive(): boolean {
    return this.connection?.status?.isActive || false;
  }

  get statusText(): string {
    return this.isActive ? 'Connected' : 'Disconnected';
  }

  get statusClass(): string {
    return this.isActive ? 'status-active' : 'status-inactive';
  }

  getClientBadgeClass(): string {
    const clientName = this.connection?.vpnClientName?.toLowerCase() || '';
    
    if (clientName.includes('wireguard') || clientName === 'wireguard') {
      return 'client-badge-wireguard';
    } else if (clientName.includes('vpnc') || clientName === 'vpnc') {
      return 'client-badge-vpnc';
    } else {
      return 'client-badge-default';
    }
  }

  onToggleConnection(): void {
    if (this.isLoading) return;
    
    if (this.isActive) {
      this.disconnect.emit(this.connection);
    } else {
      this.connect.emit(this.connection);
    }
  }
} 