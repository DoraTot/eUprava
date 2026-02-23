import { Component, OnInit } from '@angular/core';
import { AuthService } from '@auth0/auth0-angular';
import { NotificationsService } from '../../notifications/notifications.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  role: string = "";
  notificationsCount = 0;

  constructor(public auth: AuthService, private notifications: NotificationsService) {}

  ngOnInit() {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.role = claims['https://myapp.example/role'];
      }
    });

    this.notifications.notifications$.subscribe(notifs => {
      this.notificationsCount = notifs.length;
    });
  }

  logout() {
    this.auth.logout({ logoutParams: { returnTo: window.location.origin } });
  }
}
