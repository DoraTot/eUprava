import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

interface Appointment {
  id?: number;
  child_name: string;
  parent_id: number;
  doctor_id: number;
  date_time: string;
  notes?: string;
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


  constructor(private fb: FormBuilder, private http: HttpClient) {}

  ngOnInit(): void {
    this.appointmentForm = this.fb.group({
      child_name: ['', Validators.required],
      parent_id: ["", Validators.required],
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

    this.loadAppointments();
    this.fetchParentsDoctors();
  }

  loadAppointments() {
    const parentId = 1;
    this.http.get<Appointment[]>(`http://localhost:8081/getAppointments`)
      .subscribe(data => this.appointments = data);
  }

  openModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'block';
    el.classList.add('show', 'd-block');
    document.body.classList.add('modal-open');
  }


  openModal2(appointment: Appointment) {
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
    this.http.post('http://localhost:8081/createAppointment', newAppointment)
      .subscribe(() => {
        this.loadAppointments();
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

    this.http.post('http://localhost:8081/createJustification', newMedicalRecord)
      .subscribe(() => {
        this.loadAppointments();
        this.closeModal();
        this.medForm.reset();
      });

  }


  fetchParentsDoctors() {
    this.http.get<any[]>('http://localhost:8082/parents').subscribe(data => {
      this.parents = data;
    });
    this.http.get<any[]>('http://localhost:8082/doctors').subscribe(data => {
      this.doctors = data;
    });

  }

}
