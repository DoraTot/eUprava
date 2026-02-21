import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {HttpClient} from '@angular/common/http';
import {AuthService} from '@auth0/auth0-angular';

interface MedicalJustification {
  id?: number;
  child_name: string;
  doctor_id: string;
  parent_id: string;
  valid_from: string;
  valid_to: string;
  reason: string;
}

@Component({
  selector: 'app-med-justification',
  templateUrl: './med-justification.component.html',
  styleUrl: './med-justification.component.css'
})
export class MedJustificationComponent implements OnInit {
  justifications: MedicalJustification[] = [];
  justificationForm!: FormGroup;
  @ViewChild('justificationModal') modal!: ElementRef;

  authIdToken: string | null = null;
  role: string = "";
  user_id: string = "";

  constructor(private fb: FormBuilder, private http: HttpClient, public auth: AuthService) {}

  ngOnInit(): void {
    this.justificationForm = this.fb.group({
      childName: ['', Validators.required],
      doctorID: [0, Validators.required],
      parentID: [0, Validators.required],
      date: ['', Validators.required],
      reason: ['', Validators.required]
    });

    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.authIdToken = claims.__raw;
        const role = claims['https://myapp.example/role'];
        console.log('Auth0 Claims:', claims);
        this.role = role;
        this.user_id = claims['sub'];
        if (this.role == "Doctor") {
          this.loadJustificationsForDoctor();
        } else if (this.role != "Educator") {
          this.loadJustificationsForParent();
        } else {
          this.loadAllJustifications();
        }
      }
    });

  }

  loadAllJustifications() {
    this.http.get<MedicalJustification[]>(`http://localhost:8081/getAllJustifications`)
      .subscribe(data => this.justifications = data);
    console.log(this.justifications);
  }

  loadJustificationsForDoctor() {
    this.http.get<MedicalJustification[]>(`http://localhost:8081/getJustificationsForDoctor/` + this.user_id)
      .subscribe(data => this.justifications = data);
    console.log(this.justifications);
  }

  loadJustificationsForParent() {
    this.http.get<MedicalJustification[]>(`http://localhost:8081/getJustificationsForParent/` + this.user_id)
      .subscribe(data => this.justifications = data);
    console.log(this.justifications);
  }

  addJustification() {
    if (this.justificationForm.invalid) return;

    const newJustification = this.justificationForm.value;
    this.http.post('http://localhost:8081/createJustification', newJustification)
      .subscribe(() => {
        this.justificationForm.reset();
      });
  }

  openModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'block';
    el.classList.add('show');
    el.classList.add('d-block'); // ensure modal is visible
    document.body.classList.add('modal-open'); // prevent background scroll
  }

  closeModal() {
    const el = this.modal.nativeElement;
    el.style.display = 'none';
    el.classList.remove('show');
    el.classList.remove('d-block');
    document.body.classList.remove('modal-open');
  }

}
