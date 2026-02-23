import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {ActivatedRoute} from '@angular/router';
import {AuthService} from '@auth0/auth0-angular';

@Component({
  selector: 'app-statistics',
  templateUrl: './statistics.component.html',
  styleUrl: './statistics.component.css'
})
export class StatisticsComponent implements OnInit {
  totalRecords = 0;
  totalAbsent = 0;
  totalPresent = 0;
  totalJustified = 0;

  monthlyAppointments = 0;
  yearlyAppointments = 0;
  authIdToken: string | null = null;
  role: string = "";
  childName = '';

  constructor(private route: ActivatedRoute, private http: HttpClient, public auth: AuthService) {}

  ngOnInit(): void {
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'];
        this.role = role;
      }
    });
    this.route.params.subscribe(params => {
      this.childName = params['childName'];
      this.getStatistics();
    });
  }

  getStatistics() {
    const url = `http://localhost:8080/getChildStats/${this.childName}`;
    this.http.get<any>(url).subscribe({
      next: (data) => {
        this.totalRecords = data.totalRecords;
        this.totalAbsent = data.totalAbsent;
        this.totalPresent = data.totalPresent;
        this.totalJustified = data.totalJustified;

        this.monthlyAppointments = data.monthlyAppointments;
        this.yearlyAppointments = data.yearlyAppointments;
      },
      error: (err) => {
        console.error('Failed to load statistics', err);
      }
    });
  }
}
