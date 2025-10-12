import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AttemdanceRecordComponent } from './attemdance-record.component';

describe('AttemdanceRecordComponent', () => {
  let component: AttemdanceRecordComponent;
  let fixture: ComponentFixture<AttemdanceRecordComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AttemdanceRecordComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AttemdanceRecordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
