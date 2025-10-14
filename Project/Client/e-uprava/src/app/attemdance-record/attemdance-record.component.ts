import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { AuthService } from '@auth0/auth0-angular';
import { HttpClient } from '@angular/common/http';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-attemdance-record',
  templateUrl: './attemdance-record.component.html',
  styleUrls: ['./attemdance-record.component.css']
})
export class AttemdanceRecordComponent implements OnInit {

  @ViewChild('attendanceModal') modal!: ElementRef;

  attendanceForm!: FormGroup;
  records: any[] = [];

  constructor(private fb: FormBuilder, public auth: AuthService, private http: HttpClient) {}

  ngOnInit(): void {
    this.attendanceForm = this.fb.group({
      child: ['', Validators.required],
      parent: ['', Validators.required],
      date: ['', Validators.required],
      missing: [false]
    });

    this.loadRecords();
  }

  loadRecords() {
    this.http.get<any[]>('http://localhost:8080/attendance')
      .subscribe(res => this.records = res);
  }

  openModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'block';
    el.classList.add('show', 'd-block');
    document.body.classList.add('modal-open');
  }

  closeModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'none';
    el.classList.remove('show', 'd-block');
    document.body.classList.remove('modal-open');
  }

  addRecord() {
    if (this.attendanceForm.invalid) return;

    const newRecord = this.attendanceForm.value;
    this.http.post('http://localhost:8080/attendance', newRecord)
      .subscribe({
        next: () => {
          this.loadRecords();
          this.attendanceForm.reset({ missing: false });
          this.closeModal();
        },
        error: err => console.error('Failed to add record', err)
      });
  }

  markMissing(id: number) {
    this.http.post(`http://localhost:8080/attendance/${id}/missing`, {})
      .subscribe(() => this.loadRecords());
  }

  justify(id: number) {
    this.http.post(`http://localhost:8080/attendance/${id}/justify`, {})
      .subscribe(() => this.loadRecords());
  }

}
