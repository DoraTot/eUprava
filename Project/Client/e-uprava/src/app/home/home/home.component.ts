import {Component, OnInit} from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent  implements OnInit {
  role: string = "";
  userId: string = "";
  authIdToken: string | null = null;

  constructor(public auth: AuthService) {
  }

  ngOnInit(){
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'];
        this.userId = claims['sub'];
        this.role = role;
      }

    });
  }
}
