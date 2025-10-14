import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginPatientComponent } from './login-patient/login-patient.component'
import {ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule} from '@angular/common/http';
import {AuthModule} from '@auth0/auth0-angular';
import { AttemdanceRecordComponent } from './attemdance-record/attemdance-record.component';
import { HomeComponent } from './home/home/home.component';
import { AppointmentsComponent } from './appointments/appointments/appointments.component';
import { MedJustificationComponent } from './med-justification/med-justification/med-justification.component';
import { HeaderComponent } from './header/header/header.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginPatientComponent,
    AttemdanceRecordComponent,
    HomeComponent,
    AppointmentsComponent,
    MedJustificationComponent,
    HeaderComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule,
    AuthModule.forRoot({
      domain: 'dev-h1l4uuvj170yqf56.us.auth0.com',
      clientId: '16C4j8pjNYynw5qTDX9LOK5LC2nDu3wx',
      authorizationParams: {
        redirect_uri: window.location.origin,
        audience: 'https://dev-h1l4uuvj170yqf56.us.auth0.com/api/v2/',
      },
      useRefreshTokens: true,
      cacheLocation: 'localstorage'
    }),

  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
