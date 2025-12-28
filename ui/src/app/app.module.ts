import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { provideZoneChangeDetection } from '@angular/core';

import { AppComponent } from './app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { ConnectionListComponent } from './components/connection-list/connection-list.component';
import { ConnectionCardComponent } from './components/connection-card/connection-card.component';
import { StatusCardComponent } from './components/status-card/status-card.component';
import { HeaderComponent } from './components/header/header.component';

import { VpnService } from './services/vpn.service';

@NgModule({ declarations: [
        AppComponent,
        DashboardComponent,
        ConnectionListComponent,
        ConnectionCardComponent,
        StatusCardComponent,
        HeaderComponent
    ],
    bootstrap: [AppComponent], imports: [BrowserModule,
        HttpClientModule,
        FormsModule,
        ReactiveFormsModule], providers: [VpnService, provideZoneChangeDetection({ eventCoalescing: true })] })
export class AppModule { } 