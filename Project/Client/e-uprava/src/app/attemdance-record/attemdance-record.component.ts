import {booleanAttribute, Component} from '@angular/core';
import {AuthService} from '@auth0/auth0-angular';
import {HttpClient} from '@angular/common/http';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';


@Component({
  selector: 'app-attemdance-record',
  templateUrl: './attemdance-record.component.html',
  styleUrl: './attemdance-record.component.css'
})
export class AttemdanceRecordComponent {

  attendanceForm: FormGroup;

  constructor(private fb: FormBuilder,public auth: AuthService, private http: HttpClient) {
    this.attendanceForm = this.fb.group({
      child: ['', Validators.required],
      parent: ['', Validators.required],
      date: ['', Validators.required]
    });
  };

  logout() {
    this.auth.logout({ logoutParams: { returnTo: window.location.origin } });
  }

  records: any[] = [];

  newRecord = {
    child: '',
    parent: '',
    date: '',
    missing: false
  };

  ngOnInit(): void {
    this.loadRecords();
    this.auth.user$.subscribe(user => {
      if (user) {
        const role = user['https://myapp.example/role'];
        console.log('User role:', role);
      }
    });

    this.attendanceForm = this.fb.group({
      child: ['', Validators.required],
      parent: ['', Validators.required],
      date: ['', Validators.required],
      missing: [false]
    });
  }

  loadRecords() {
    this.http.get<any[]>('http://localhost:8080/attendance')
      .subscribe(res => this.records = res);
  }

  addRecord() {
    this.http.post('http://localhost:8080/attendance', this.newRecord)
      .subscribe({
        next: () => {
          this.loadRecords();
          // this.newRecord = {child: '', parent: '', date: '', missing: false};
          this.attendanceForm.reset({ missing: false });
        },
        error: err => console.error('Failed to add record', err)
      });
  }

  markMissing(id: number) {
    this.http.post(`http://localhost:8080/attendance/${id}/missing`, {})
      .subscribe(() => this.loadRecords()); // refresh after update
  }

  justify(id: number) {
    this.http.post(`http://localhost:8080/attendance/${id}/justify`, {})
      .subscribe(() => this.loadRecords()); // refresh after update
  }



}
