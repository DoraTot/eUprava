import {Component, OnInit} from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent implements OnInit {
  role: string = "";
  constructor(public auth: AuthService) {
  }

  ngOnInit() {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.role = claims['https://myapp.example/role'];
      }
    });
  }
  logout() {
    this.auth.logout({ logoutParams: { returnTo: window.location.origin } });
  }
}
