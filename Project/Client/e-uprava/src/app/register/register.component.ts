import {Component, OnInit} from '@angular/core';
import {User} from "../model/user.model";
import {HttpClient} from '@angular/common/http';
import {FormGroup, FormBuilder, Validators, ReactiveFormsModule} from '@angular/forms';
import {Router} from '@angular/router';


@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent implements OnInit{

  registrationForm!: FormGroup;
  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {}

  ngOnInit() {
    this.registrationForm = this.fb.group({
      username: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required],
      fname: ['', Validators.required],
      lname: ['', Validators.required],
      usertype: ['']
    });
  }


  onSubmit() {
    if (this.registrationForm.valid) {
      const formData = {
        username: this.registrationForm.value.username,
        password: this.registrationForm.value.password,
        fname: this.registrationForm.value.fname,
        lname: this.registrationForm.value.lname,
        usertype: this.registrationForm.value.usertype
      };
      console.log('Sending registration data:', formData);
      this.http.post('http://localhost:8080/register', formData).subscribe({
        next: () => {
          console.log('User registered successfully!');
          this.router.navigate(['/login-patient']);
        },
        error: () => alert('Registration failed')
      });
    }
  }

}
