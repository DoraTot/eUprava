import {Component, OnDestroy, OnInit} from '@angular/core';
import { NotificationsService } from './notifications.service';
import { Observable } from 'rxjs';
import {AuthService} from '@auth0/auth0-angular';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-notifications',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.css']
})
export class NotificationsComponent implements OnInit, OnDestroy {
  notifications$: Observable<string[]>;
  authIdToken: string | null = null;
  user_id: string = "";
  oldNotifications: any[] = [];

  constructor(private notifications: NotificationsService, private auth: AuthService, private http: HttpClient) {
    this.notifications$ = this.notifications.notifications$;
  }

  ngOnInit(): void {

    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        this.user_id = claims['sub'];

        this.http.get<any[]>(`http://localhost:8080/notifications/${this.user_id}`)
          .subscribe(data => {
            this.oldNotifications = data;
            this.markAllAsRead();
          });
      }
    });
  }

  ngOnDestroy(){
    this.notifications.clearAll();
  }

  markAllAsRead() {
    this.http.put(`http://localhost:8080/notifications/read/${this.user_id}`, {})
      .subscribe({
        error: err => console.error("Failed to mark as read", err)
      });

  }


}
