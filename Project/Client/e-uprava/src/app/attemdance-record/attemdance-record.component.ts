import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { AuthService } from '@auth0/auth0-angular';
import { HttpClient } from '@angular/common/http';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';


interface Attendance {
  id?: number;
  child_name: string;
  parent_id: number;
  doctor_id: number;
  date: string;
  medical_record_id: string;
}

@Component({
  selector: 'app-attemdance-record',
  templateUrl: './attemdance-record.component.html',
  styleUrls: ['./attemdance-record.component.css']
})


export class AttemdanceRecordComponent implements OnInit {

  @ViewChild('attendanceModal') modal!: ElementRef;
  @ViewChild('appointmentModal2') modal2!: ElementRef;

  attendanceForm!: FormGroup;
  appForm: FormGroup = new FormGroup({});
  medicalRecords: [] = [];
  viewMode: 'today' | 'history' = 'today';
  today = new Date();
  records: any[] = [];
  parents: any[] = [];
  authIdToken: string | null = null;
  role: string = "";
  children: any[] = [];

  userId: string = "";

  constructor(private fb: FormBuilder, public auth: AuthService, private http: HttpClient) {}

  ngOnInit(): void {
    this.attendanceForm = this.fb.group({
      child: ['', Validators.required],
      missing: [false]
    });
    this.appForm = this.fb.group({
      medical_justification: ['', Validators.required]
    });


    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'];
        console.log('Auth0 ID Token:', this.authIdToken);
        console.log('Auth0 Claims:', claims);
        console.log('User role:', role);
        this.userId = claims['sub'];
        this.role = role;
        this.fetchChildren();


        if (this.viewMode === 'today') {
          if (this.role == "Educator") {
            this.loadRecordsForToday();
          } else if (this.role != "Doctor") {
            this.loadRecordsByParentForToday();
          }
        } else {
          if (this.role == "Educator") {
            this.loadRecords();
          } else if (this.role != "Doctor") {
            this.loadRecordsByParent();
          }
        }

      }

    });


    this.fetchParents();
    // this.loadMedicalRecords(this.userId)
  }

  setView(mode: 'today' | 'history') {
    this.viewMode = mode;

    if (mode === 'today') {
      if (this.role == "Educator") {
        this.loadRecordsForToday();
      } else if (this.role != "Doctor") {
        this.loadRecordsByParentForToday();
      }
    } else {
      if (this.role == "Educator") {
        this.loadRecords();
      } else if (this.role != "Doctor") {
        this.loadRecordsByParent();
      }
    }
  }

  fetchChildren() {
    this.http.get<any[]>(`http://localhost:8080/children/getAll`).subscribe(data => {
      this.children = data;
    });
  }
  loadMedicalRecords(userId: string) {

    // this.auth.idTokenClaims$.subscribe(claims => {
    //   if (claims && claims.__raw) {
    //     this.authIdToken = claims.__raw;
    //     const role = claims['https://myapp.example/role'];
    //     const userSub = claims['sub'] as string;
    //     console.log('Auth0 ID Token:', this.authIdToken);
    //     console.log('Auth0 Claims:', claims);
    //     console.log('User role:', role);
    //     this.role = role;
    //     // this.userId = claims.sub;
    //     this.http.get<any[]>(`http://localhost:8081/medicalRecord/user/${userSub}`)
    //       .subscribe(res => this.records = res,
    //         console.log(this.records););
    //
    //   }
    // });
    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'] as string;
        const userSub = claims['sub'] as string;
        this.role = role;

        this.http
          .get<any[]>(`http://localhost:8081/medicalRecord/user/${userSub}`)
          .subscribe(res => {
            this.records = res;
            console.log("records: ", this.records);
          });
      }
    });

  }

  loadRecordsForToday() {
    this.http.get<any[]>('http://localhost:8080/attendanceForToday')
      .subscribe(res => this.records = res);
  }

  loadRecords() {
    this.http.get<any[]>('http://localhost:8080/attendance')
      .subscribe(res => this.records = res);
  }

  loadRecordsByParent() {
    this.http.get<any[]>('http://localhost:8080/attendanceByParent/' + this.userId)
      .subscribe(res => this.records = res);
  }

  loadRecordsByParentForToday() {
    this.http.get<any[]>('http://localhost:8080/attendanceByParentForToday/' + this.userId)
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
    newRecord.date = (new Date()).toISOString().slice(0, 10);

    const selectedChildName = this.attendanceForm.value.child;

    const selectedChild = this.children.find(
      c => c.name === selectedChildName
    );
    newRecord.parent = selectedChild.parent_id;
    newRecord.justified = false;
    console.log(newRecord);
    this.http.post('http://localhost:8080/attendance', newRecord)
      .subscribe({
        next: () => {
          if (this.viewMode === 'today') {
            if (this.role == "Educator") {
              this.loadRecordsForToday();
            } else if (this.role != "Doctor") {
              this.loadRecordsByParentForToday();
            }
          } else {
            if (this.role == "Educator") {
              this.loadRecords();
            } else if (this.role != "Doctor") {
              this.loadRecordsByParent();
            }
          }
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

  justify(id: string) {
    this.http.post(`http://localhost:8080/attendance/${id}/justify`, {})
      .subscribe(() => this.loadRecords());
  }

  fetchParents() {
    this.http.get<any[]>('http://localhost:8082/parents').subscribe(data => {
      this.parents = data;
    });
  }

  openModal2(attendanceId: string) {
    const el = this.modal2.nativeElement;
    el.style.display = 'block';
    el.classList.add('show', 'd-block');
    document.body.classList.add('modal-open');
    console.log("userID: " + this.userId);


    this.loadMedicalRecords(this.userId);
  }

  closeModal2() {
    const el = this.modal2.nativeElement;
    el.style.display = 'none';
    el.classList.remove('show', 'd-block');
    document.body.classList.remove('modal-open');
  }

}
