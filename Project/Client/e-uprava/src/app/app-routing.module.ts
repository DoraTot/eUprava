import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginPatientComponent } from './login-patient/login-patient.component';
import { LoginDoctorComponent } from './login-doctor/login-doctor.component';

// const routes: Routes = [];

//
// export const routes: Routes = [
//   { path: 'login-patient', component: RegistrationComponent, canActivate:[loginGuard]},
//   { path: '', redirectTo: '/register', pathMatch: 'full' }
//
// ];

export const routes: Routes = [
  { path: 'login-patient', component: LoginPatientComponent },
  { path: 'login-doctor', component: LoginDoctorComponent },

  { path: '', redirectTo: '/login-patient', pathMatch: 'full' },
];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
