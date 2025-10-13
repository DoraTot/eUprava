import {Component, OnInit} from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';
import {HttpClient} from '@angular/common/http';
import {Router} from '@angular/router';

@Component({
  selector: 'app-login-patient',
  templateUrl: './login-patient.component.html',
  styleUrl: './login-patient.component.css'
})
export class LoginPatientComponent implements OnInit{

  constructor(public auth: AuthService, private http: HttpClient, private router: Router) {}

  selectedRole: 'Patient' | 'Doctor' | null = null;
  authIdToken: string | null = null;

  ngOnInit() {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'];
        console.log('Auth0 ID Token:', this.authIdToken);
        console.log('Auth0 Claims:', claims);
        console.log('User role:', role);
        // if (role === 'Educator') {
          this.router.navigate(['/attendance']);
        // }
        // else {
        //   this.router.navigate(['/attendance']);
        // }
      }
    });
  }

  login() {
    this.auth.loginWithRedirect();
  }

  logout() {
    this.auth.logout({ logoutParams: { returnTo: window.location.origin } });
  }


}
