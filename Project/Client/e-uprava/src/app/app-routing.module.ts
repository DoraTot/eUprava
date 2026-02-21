import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginPatientComponent } from './login-patient/login-patient.component';
import {AttemdanceRecordComponent} from './attemdance-record/attemdance-record.component';
import {HomeComponent} from './home/home/home.component';
import {AppointmentsComponent} from './appointments/appointments/appointments.component';
import {MedJustificationComponent} from './med-justification/med-justification/med-justification.component';
import {EnrollmentComponent} from './enrollment/enrollment.component';


export const routes: Routes = [
  { path: 'login', component: LoginPatientComponent },
  { path: 'attendance', component: AttemdanceRecordComponent },
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'appointments', component: AppointmentsComponent },
  { path: 'medicalJustification', component: MedJustificationComponent },
  { path: 'enrollment', component: EnrollmentComponent },

];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
