import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

interface Appointment {
  id?: number;
  childName: string;
  parentID: number;
  doctorID: number;
  dateTime: string;
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

  @ViewChild('appointmentModal') modal!: ElementRef;

  constructor(private fb: FormBuilder, private http: HttpClient) {}

  ngOnInit(): void {
    this.appointmentForm = this.fb.group({
      childName: ['', Validators.required],
      parentID: [0, Validators.required],
      doctorID: [0, Validators.required],
      dateTime: ['', Validators.required],
      notes: ['']
    });

    this.loadAppointments();
  }

  loadAppointments() {
    const parentId = 1; // replace with real logged-in parent ID
    this.http.get<Appointment[]>(`http://localhost:8081/getAppointments?parentId=${parentId}`)
      .subscribe(data => this.appointments = data);
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

  addAppointment() {
    if (this.appointmentForm.invalid) return;

    const newAppointment = this.appointmentForm.value;
    this.http.post('http://localhost:8081/createAppointment', newAppointment)
      .subscribe(() => {
        this.loadAppointments();
        this.closeModal();
        this.appointmentForm.reset();
      });
  }
}
