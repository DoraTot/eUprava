import {Component, OnInit} from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';
import {NavigationEnd, Router} from '@angular/router';
import {filter} from 'rxjs';
import {NotificationsService} from './notifications/notifications.service';

@Component({
  selector: 'app-root',
  template: `
    <app-header *ngIf="showHeader"></app-header>
    <router-outlet></router-outlet>
  `,
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'e-uprava';
  showHeader = true;
  authIdToken: string | null = null;
  user_id: string = "";

  constructor(private router: Router, private auth: AuthService, private notifications: NotificationsService) {
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe((event: NavigationEnd) => {
      this.showHeader = event.urlAfterRedirects !== '/login';
    });
  }

  ngOnInit(): void {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        this.user_id = claims['sub'];

        const wsPreschool = new WebSocket(`ws://localhost:8080/ws?userId=${this.user_id}`);
        wsPreschool.onmessage = (event) => {
          this.notifications.addNotification(`${event.data}`);
        };

        const wsHealthcare = new WebSocket(`ws://localhost:8081/ws?userId=${this.user_id}`);
        wsHealthcare.onmessage = (event) => {
          this.notifications.addNotification(`${event.data}`);
        };
      }
    });
  }
}
