import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationsService {
  private notificationsSubject = new BehaviorSubject<string[]>([]);
  notifications$ = this.notificationsSubject.asObservable();

  addNotification(message: string) {
    const current = this.notificationsSubject.value;
    this.notificationsSubject.next([message, ...current]);
  }

  clearAll() {
    this.notificationsSubject.next([]);
  }

  getCount(): number {
    return this.notificationsSubject.value.length;
  }
}
