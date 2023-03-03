import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentstatusComponent } from './torrentstatus.component';

describe('TorrentstatusComponent', () => {
  let component: TorrentstatusComponent;
  let fixture: ComponentFixture<TorrentstatusComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TorrentstatusComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TorrentstatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
