import { Component, Input, Output, EventEmitter } from '@angular/core';
import { VPNConfig } from '../../models/api.models';

@Component({
  selector: 'app-connection-list',
  templateUrl: './connection-list.component.html',
  styleUrls: ['./connection-list.component.scss']
})
export class ConnectionListComponent {
  @Input() connections: VPNConfig[] = [];
  @Input() isLoading: boolean = false;
  @Output() connect = new EventEmitter<VPNConfig>();
  @Output() disconnect = new EventEmitter<VPNConfig>();
  @Output() refresh = new EventEmitter<void>();

  onConnect(connection: VPNConfig): void {
    this.connect.emit(connection);
  }

  onDisconnect(connection: VPNConfig): void {
    this.disconnect.emit(connection);
  }

  onRefresh(): void {
    this.refresh.emit();
  }

  trackByConnectionId(index: number, connection: VPNConfig): string {
    return connection.id;
  }
} 