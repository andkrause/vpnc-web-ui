import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Status, VPNConfig, ConnectionStatus, DesiredConnectionStatus, ApiError } from '../models/api.models';

@Injectable({
  providedIn: 'root'
})
export class VpnService {
  private readonly apiBaseUrl = '/api/v1';
  private statusSubject = new BehaviorSubject<Status | null>(null);
  private connectionsSubject = new BehaviorSubject<VPNConfig[]>([]);

  public status$ = this.statusSubject.asObservable();
  public connections$ = this.connectionsSubject.asObservable();

  constructor(private http: HttpClient) { }

  /**
   * Get overall status of the VPN Gateway
   */
  getOverallStatus(): Observable<Status> {
    return this.http.get<Status>(`${this.apiBaseUrl}/`)
      .pipe(
        tap(status => this.statusSubject.next(status)),
        catchError(this.handleError)
      );
  }

  /**
   * Get list of all VPN connections
   */
  getConnections(): Observable<VPNConfig[]> {
    return this.http.get<VPNConfig[]>(`${this.apiBaseUrl}/connections`)
      .pipe(
        tap(connections => this.connectionsSubject.next(connections)),
        catchError(this.handleError)
      );
  }

  /**
   * Get status of a specific connection
   */
  getConnectionStatus(client: string, id: string): Observable<ConnectionStatus> {
    return this.http.get<ConnectionStatus>(`${this.apiBaseUrl}/connections/connection/${client}/${id}/`)
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Set the desired state of a connection
   */
  setConnectionStatus(client: string, id: string, desiredStatus: DesiredConnectionStatus): Observable<ConnectionStatus> {
    return this.http.post<ConnectionStatus>(
      `${this.apiBaseUrl}/connections/connection/${client}/${id}/`,
      desiredStatus
    ).pipe(
      tap(() => {
        // Refresh connections after status change
        this.getConnections().subscribe();
      }),
      catchError(this.handleError)
    );
  }

  /**
   * Connect to a VPN
   */
  connectVpn(client: string, id: string): Observable<ConnectionStatus> {
    return this.setConnectionStatus(client, id, { desiredConnectionState: 'active' });
  }

  /**
   * Disconnect from a VPN
   */
  disconnectVpn(client: string, id: string): Observable<ConnectionStatus> {
    return this.setConnectionStatus(client, id, { desiredConnectionState: 'inactive' });
  }

  /**
   * Refresh all data
   */
  refreshData(): void {
    this.getOverallStatus().subscribe();
    this.getConnections().subscribe();
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'An unknown error occurred';
    
    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Error: ${error.error.message}`;
    } else {
      // Server-side error
      if (error.error && error.error.message) {
        errorMessage = error.error.message;
      } else {
        errorMessage = `Error Code: ${error.status}\nMessage: ${error.message}`;
      }
    }
    
    console.error('VPN Service Error:', errorMessage);
    return throwError(() => new Error(errorMessage));
  }
} 