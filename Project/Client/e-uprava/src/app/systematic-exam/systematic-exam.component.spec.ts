import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SystematicExamComponent } from './systematic-exam.component';

describe('SystematicExamComponent', () => {
  let component: SystematicExamComponent;
  let fixture: ComponentFixture<SystematicExamComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystematicExamComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SystematicExamComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
