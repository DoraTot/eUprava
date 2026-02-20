import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import {AuthService} from '@auth0/auth0-angular';
import {Router} from '@angular/router';

interface Appointment {
  id?: number;
  child_name: string;
  parent_id: number;
  doctor_id: number;
  date_time: string;
  notes?: string;
  justified: boolean;
}

@Component({
  selector: 'app-appointments',
  templateUrl: './appointments.component.html',
  styleUrls: ['./appointments.component.css']
})
export class AppointmentsComponent implements OnInit {
  appointments: Appointment[] = [];
  appointmentForm!: FormGroup;
    medForm: FormGroup = new FormGroup({});
  parents: any[] = [];
  doctors: any[] = [];

  @ViewChild('appointmentModal') modal!: ElementRef;
  @ViewChild('appointmentModal2') modal2!: ElementRef;


  constructor(private fb: FormBuilder, private http: HttpClient, private auth: AuthService) {}
  authIdToken: string | null = null;
  role: string = "";
  user_id: string = "";
  minDateTime: string = '';
  selectedAppointmentId!: number;

  ngOnInit(): void {

    const now = new Date();
    const year = now.getFullYear();
    const month = (now.getMonth() + 1).toString().padStart(2, '0');
    const day = now.getDate().toString().padStart(2, '0');
    const hours = now.getHours().toString().padStart(2, '0');
    const minutes = now.getMinutes().toString().padStart(2, '0');

    this.minDateTime = `${year}-${month}-${day}T${hours}:${minutes}`;

    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        this.role = claims['https://myapp.example/role'];
        this.user_id = claims['sub'];
        if (this.role == "Doctor") {
          this.loadAppointmentsByDoctor();
        } else if (this.role != "Educator") {
          this.loadAppointmentsByParent();
        }
      }
    });

    this.appointmentForm = this.fb.group({
      child_name: ['', Validators.required],
      doctor_id: [0, Validators.required],
      date_time: ['', Validators.required],
      notes: ['']
    });

    this.medForm = this.fb.group({
      child_name: ['', Validators.required],
      parent_id: ["", Validators.required],
      doctor_id: ["", Validators.required],
      dated: ['', Validators.required],
      reason: ['']
    });

    this.fetchParentsDoctors();
  }

  loadAppointmentsByParent() {
    this.http.get<Appointment[]>(`http://localhost:8081/getAppointmentsByParent/` + this.user_id)
      .subscribe(data => this.appointments = data);
  }

  loadAppointmentsByDoctor() {
    this.http.get<Appointment[]>(`http://localhost:8081/getAppointmentsByDoctor/` + this.user_id)
      .subscribe(data => this.appointments = data);
  }

  openModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'block';
    el.classList.add('show', 'd-block');
    document.body.classList.add('modal-open');
  }


  openModal2(appointment: Appointment) {
    this.selectedAppointmentId = appointment.id!;

    const el = this.modal2.nativeElement;
    el.style.display = 'block';
    el.classList.add('show', 'd-block');
    document.body.classList.add('modal-open');

    this.medForm.patchValue({
      child_name: appointment.child_name,
      parent_id: appointment.parent_id,
      doctor_id: appointment.doctor_id,
      dated: appointment.date_time,
      reason: '',
    });
  }

  closeModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'none';
    el.classList.remove('show', 'd-block');
    document.body.classList.remove('modal-open');
  }

  closeModal2() {
    const el = this.modal2.nativeElement;
    el.style.display = 'none';
    el.classList.remove('show', 'd-block');
    document.body.classList.remove('modal-open');
  }

  addAppointment() {
    console.log(this.appointmentForm.value);
    if (this.appointmentForm.invalid) return;

    const newAppointment = this.appointmentForm.value;
    console.log(newAppointment);
    newAppointment.parent_id = this.user_id;
    newAppointment.justified = false;
    this.http.post('http://localhost:8081/createAppointment', newAppointment)
      .subscribe(() => {
        if (this.role == "Doctor") {
          this.loadAppointmentsByDoctor();
        } else {
          this.loadAppointmentsByParent();
        }
        this.closeModal();
        this.appointmentForm.reset();
      });
  }

  approve(){
    // console.log('Form values:', this.medForm.value);
    const newMedicalRecord = this.medForm.value;
    if (newMedicalRecord.dated) {
      const date = new Date(newMedicalRecord.dated);
      const year = date.getFullYear();
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const day = date.getDate().toString().padStart(2, '0');
      newMedicalRecord.dated = `${year}-${month}-${day}`;
    }

    console.log('Form values:', newMedicalRecord);

    // this.http.post('http://localhost:8081/createJustification', newMedicalRecord)
    //   .subscribe(() => {
    //     this.loadAppointments();
    //     this.closeModal();
    //     this.medForm.reset();
    //   });

    this.http.post('http://localhost:8081/createJustification', newMedicalRecord)
      .subscribe({
        next: () => {

          this.justifyAppointment(this.selectedAppointmentId)
            .subscribe(() => {

              if (this.role == "Doctor") {
                this.loadAppointmentsByDoctor();
              } else {
                this.loadAppointmentsByParent();
              }

              this.closeModal2();
              this.medForm.reset();
            });

        },
        error: (err) => {
          console.error('POST ERROR:', err);
        }
      });
  }

  justifyAppointment(id: number) {
    return this.http.put(`http://localhost:8081/justifyAppointment/${id}`, {});
  }

  fetchParentsDoctors() {
    this.http.get<any[]>('http://localhost:8082/parents').subscribe(data => {
      this.parents = data;
    });
    this.http.get<any[]>('http://localhost:8082/doctors').subscribe(data => {
      this.doctors = data;
    });

  }

  cancelAppointment(appointmentId: any) {
    this.http.delete(`http://localhost:8081/cancelAppointment/${appointmentId}`)
      .subscribe(() => {
        if (this.role === 'Doctor') {
          this.loadAppointmentsByDoctor();
        } else {
          this.loadAppointmentsByParent();
        }
      });
  }

  isFutureDate(dateStr: string): boolean {
    const date = new Date(dateStr);
    const now = new Date();
    return date > now;
  }

}
