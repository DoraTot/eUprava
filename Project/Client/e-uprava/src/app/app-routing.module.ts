import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginPatientComponent } from './login-patient/login-patient.component';
import {AttemdanceRecordComponent} from './attemdance-record/attemdance-record.component';


export const routes: Routes = [
  { path: 'login', component: LoginPatientComponent },
  { path: 'attendance', component: AttemdanceRecordComponent },
  { path: '', redirectTo: '/login', pathMatch: 'full' },
];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
