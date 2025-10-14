import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MedJustificationComponent } from './med-justification.component';

describe('MedJustificationComponent', () => {
  let component: MedJustificationComponent;
  let fixture: ComponentFixture<MedJustificationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [MedJustificationComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MedJustificationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
