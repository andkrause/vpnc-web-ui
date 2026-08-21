import { Component, OnInit, OnDestroy, signal } from '@angular/core';
import { Subject, takeUntil, finalize } from 'rxjs';
import { VpnService } from '../../services/vpn.service';
import { Status, VPNConfig } from '../../models/api.models';
import { StatusCardComponent } from '../status-card/status-card.component';
import { ConnectionListComponent } from '../connection-list/connection-list.component';

@Component({
    selector: 'app-dashboard',
    imports: [StatusCardComponent, ConnectionListComponent],
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject<void>();

  status = signal<Status | null>(null);
  connections = signal<VPNConfig[]>([]);
  isLoadingStatus = signal(false);
  isLoadingConnections = signal(false);
  error = signal<string | null>(null);

  constructor(private vpnService: VpnService) {}

  ngOnInit(): void {
    this.loadData();

    // Subscribe to status updates
    this.vpnService.status$
      .pipe(takeUntil(this.destroy$))
      .subscribe(status => {
        this.status.set(status);
      });

    // Subscribe to connections updates
    this.vpnService.connections$
      .pipe(takeUntil(this.destroy$))
      .subscribe(connections => {
        this.connections.set(connections);
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadData(): void {
    this.loadStatus();
    this.loadConnections();
  }

  loadStatus(): void {
    this.isLoadingStatus.set(true);
    this.error.set(null);

    this.vpnService.getOverallStatus()
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.isLoadingStatus.set(false))
      )
      .subscribe({
        next: (status) => {
          this.status.set(status);
        },
        error: (error) => {
          this.error.set(`Failed to load status: ${error.message}`);
          console.error('Failed to load status:', error);
        }
      });
  }

  loadConnections(): void {
    this.isLoadingConnections.set(true);

    this.vpnService.getConnections()
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.isLoadingConnections.set(false))
      )
      .subscribe({
        next: (connections) => {
          this.connections.set(connections);
        },
        error: (error) => {
          this.error.set(`Failed to load connections: ${error.message}`);
          console.error('Failed to load connections:', error);
        }
      });
  }

  onConnect(connection: VPNConfig): void {
    this.vpnService.connectVpn(connection.vpnClientName, connection.configName)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.loadData(); // Refresh all data
        },
        error: (error) => {
          this.error.set(`Failed to connect: ${error.message}`);
          console.error('Failed to connect:', error);
        }
      });
  }

  onDisconnect(connection: VPNConfig): void {
    this.vpnService.disconnectVpn(connection.vpnClientName, connection.configName)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.loadData(); // Refresh all data
        },
        error: (error) => {
          this.error.set(`Failed to disconnect: ${error.message}`);
          console.error('Failed to disconnect:', error);
        }
      });
  }

  onDisconnectAll(): void {
    const status = this.status();
    if (!status?.activeVpnClient || !status?.activeVpnConfig) {
      return;
    }

    this.vpnService.disconnectVpn(status.activeVpnClient, status.activeVpnConfig)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.loadData(); // Refresh all data
        },
        error: (error) => {
          this.error.set(`Failed to disconnect: ${error.message}`);
          console.error('Failed to disconnect:', error);
        }
      });
  }

  onRefresh(): void {
    this.loadData();
  }

  dismissError(): void {
    this.error.set(null);
  }
}
