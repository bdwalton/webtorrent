import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ServerconfigComponent } from './serverconfig.component';

describe('ServerconfigComponent', () => {
  let component: ServerconfigComponent;
  let fixture: ComponentFixture<ServerconfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ServerconfigComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ServerconfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
