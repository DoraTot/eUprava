import { Component } from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';

@Component({
  selector: 'app-attemdance-record',
  templateUrl: './attemdance-record.component.html',
  styleUrl: './attemdance-record.component.css'
})
export class AttemdanceRecordComponent {

  constructor(public auth: AuthService) {};

  logout() {
    this.auth.logout({ logoutParams: { returnTo: window.location.origin } });
  }

}
