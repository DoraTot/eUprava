import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {HttpClient} from '@angular/common/http';

interface MedicalJustification {
  id?: number;
  childName: string;
  doctorID: number;
  parentID: number;
  date: string;
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

  constructor(private fb: FormBuilder, private http: HttpClient) {}

  ngOnInit(): void {
    this.justificationForm = this.fb.group({
      childName: ['', Validators.required],
      doctorID: [0, Validators.required],
      parentID: [0, Validators.required],
      date: ['', Validators.required],
      reason: ['', Validators.required]
    });

    this.loadJustifications();
  }

  loadJustifications() {
    const parentId = 1; // replace with real logged-in parent ID
    this.http.get<MedicalJustification[]>(`http://localhost:8081/getJustification?parentId=${parentId}`)
      .subscribe(data => this.justifications = data);
  }

  addJustification() {
    if (this.justificationForm.invalid) return;

    const newJustification = this.justificationForm.value;
    this.http.post('http://localhost:8081/createJustification', newJustification)
      .subscribe(() => {
        this.loadJustifications();
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
