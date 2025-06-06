import { Component, Input, Output, EventEmitter } from '@angular/core';
import { Status } from '../../models/api.models';

@Component({
    selector: 'app-status-card',
    templateUrl: './status-card.component.html',
    styleUrls: ['./status-card.component.scss'],
    standalone: false
})
export class StatusCardComponent {
  @Input() status: Status | null = null;
  @Input() isLoading: boolean = false;
  @Output() disconnectAll = new EventEmitter<void>();
  @Output() refresh = new EventEmitter<void>();

  get isConnected(): boolean {
    return !!(this.status?.activeVpnClient && this.status?.activeVpnConfig);
  }

  get statusText(): string {
    if (this.isLoading) return 'Loading...';
    if (!this.status) return 'No data';
    if (this.isConnected) return 'Connected';
    return 'Disconnected';
  }

  get statusClass(): string {
    if (this.isLoading) return 'status-loading';
    if (this.isConnected) return 'status-connected';
    return 'status-disconnected';
  }

  getActiveClientBadgeClass(): string {
    const clientName = this.status?.activeVpnClient?.toLowerCase() || '';
    
    if (clientName.includes('wireguard') || clientName === 'wireguard') {
      return 'client-badge-wireguard';
    } else if (clientName.includes('vpnc') || clientName === 'vpnc') {
      return 'client-badge-vpnc';
    } else {
      return 'client-badge-default';
    }
  }

  onDisconnectAll(): void {
    this.disconnectAll.emit();
  }

  onRefresh(): void {
    this.refresh.emit();
  }
} 