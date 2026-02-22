import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {HttpClient} from '@angular/common/http';
import {AuthService} from '@auth0/auth0-angular';

interface Child {
  id?: number;
  name: string;
  parent_id: string;
  exam_done: boolean;
  enrolled: string;
}
declare var bootstrap: any;

@Component({
  selector: 'app-enrollment',
  templateUrl: './enrollment.component.html',
  styleUrl: './enrollment.component.css'
})
export class EnrollmentComponent implements OnInit {
  @ViewChild('addChildModal') addChildModal!: ElementRef;

  children: Child[] = [];
  enrollmentForm!: FormGroup;
  user_id: string = '';
  role: string = "";
  toastMessage: string = '';
  showToast: boolean = false;

  constructor(
    private fb: FormBuilder,
    private http: HttpClient,
    private auth: AuthService
  ) {}

  ngOnInit(): void {
    this.enrollmentForm = this.fb.group({
      name: ['', Validators.required]
    });

    this.auth.idTokenClaims$.subscribe(claims => {
      if (claims && claims.__raw) {
        this.user_id = claims['sub'];
        this.role = claims['https://myapp.example/role'];
        if (this.role == "Educator") {
          this.loadChildren();
        } else if (this.role != "Doctor") {
          this.loadChildrenByParent();
        }
      }
    });
  }

  loadChildrenByParent() {
    this.http.get<Child[]>(`http://localhost:8080/children/getByParentID/` + this.user_id)
      .subscribe({
        next: data => this.children = data,
        error: err => console.error('Failed to load children:', err)
      });
  }

  loadChildren() {
    this.http.get<Child[]>(`http://localhost:8080/children/getAll`)
      .subscribe({
        next: data => this.children = data,
        error: err => console.error('Failed to load children:', err)
      });
  }

  addChild() {
    if (this.enrollmentForm.invalid) return;

    const newChild = {
      name: this.enrollmentForm.value.name,
      parent_id: this.user_id
    };

    this.http.post<any>('http://localhost:8080/children/add', newChild)
      .subscribe({
        next: (res) => {
          console.log('Child created:', res);
          this.children.push({
            id: res.id,
            name: newChild.name,
            parent_id: this.user_id,
            exam_done: res.exam_done,
            enrolled: 'PENDING'
          });

          const modalEl = this.addChildModal.nativeElement;
          const modalInstance = bootstrap.Modal.getInstance(modalEl) || new bootstrap.Modal(modalEl);
          modalInstance.hide();

          this.enrollmentForm.reset();
        },
        error: err => console.error('Failed to create child:', err)
      });
  }

  enrollChild(childName: string) {

    const payload = {
      parent_id: this.user_id,
      child_name: childName
    };

    this.http.post(`http://localhost:8080/enrollChild`, payload)
      .subscribe({
        next: (res: any) => {

          if (res.exam_done) {
            this.showAutoToast(`✅ ${childName} successfully enrolled!`, 3000);
          } else {
            this.showAutoToast(`⚠️ Systematic exam not completed yet.`, 3000);
          }

          this.loadChildrenByParent();
        },
        error: (err) => {
          this.showAutoToast("❌ Enrollment failed.", 3000);
          console.error("Failed to enroll child:", err);
        }
      });
  }

  private showAutoToast(message: string, duration: number = 3000) {
    this.toastMessage = message;
    this.showToast = true;

    setTimeout(() => {
      this.showToast = false;
    }, duration);
  }

}
