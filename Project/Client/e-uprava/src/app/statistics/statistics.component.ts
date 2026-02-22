import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {ActivatedRoute} from '@angular/router';

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

  childName = '';

  constructor(private route: ActivatedRoute, private http: HttpClient) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.childName = params['childName'];
      this.getStatistics();
    });  }

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
