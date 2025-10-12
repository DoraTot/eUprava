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

  ngOnInit() {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        const token = claims.__raw;
        this.http.post('http://localhost:8080/login', { token })
          .subscribe({
            next: res => {console.log('Backend verified:', res),
              this.router.navigate(['/attendance']);},
            error: err => console.error('Verification failed:', err)
          });
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
