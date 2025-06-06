import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject, takeUntil, finalize } from 'rxjs';
import { VpnService } from '../../services/vpn.service';
import { Status, VPNConfig } from '../../models/api.models';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject<void>();
  
  status: Status | null = null;
  connections: VPNConfig[] = [];
  isLoadingStatus = false;
  isLoadingConnections = false;
  error: string | null = null;

  constructor(private vpnService: VpnService) {}

  ngOnInit(): void {
    this.loadData();
    
    // Subscribe to status updates
    this.vpnService.status$
      .pipe(takeUntil(this.destroy$))
      .subscribe(status => {
        this.status = status;
      });

    // Subscribe to connections updates
    this.vpnService.connections$
      .pipe(takeUntil(this.destroy$))
      .subscribe(connections => {
        this.connections = connections;
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
    this.isLoadingStatus = true;
    this.error = null;
    
    this.vpnService.getOverallStatus()
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.isLoadingStatus = false)
      )
      .subscribe({
        next: (status) => {
          this.status = status;
        },
        error: (error) => {
          this.error = `Failed to load status: ${error.message}`;
          console.error('Failed to load status:', error);
        }
      });
  }

  loadConnections(): void {
    this.isLoadingConnections = true;
    
    this.vpnService.getConnections()
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.isLoadingConnections = false)
      )
      .subscribe({
        next: (connections) => {
          this.connections = connections;
        },
        error: (error) => {
          this.error = `Failed to load connections: ${error.message}`;
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
          this.error = `Failed to connect: ${error.message}`;
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
          this.error = `Failed to disconnect: ${error.message}`;
          console.error('Failed to disconnect:', error);
        }
      });
  }

  onDisconnectAll(): void {
    if (!this.status?.activeVpnClient || !this.status?.activeVpnConfig) {
      return;
    }

    this.vpnService.disconnectVpn(this.status.activeVpnClient, this.status.activeVpnConfig)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.loadData(); // Refresh all data
        },
        error: (error) => {
          this.error = `Failed to disconnect: ${error.message}`;
          console.error('Failed to disconnect:', error);
        }
      });
  }

  onRefresh(): void {
    this.loadData();
  }

  dismissError(): void {
    this.error = null;
  }
} 