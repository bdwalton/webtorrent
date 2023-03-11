import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentControlsComponent } from './torrent-controls.component';

describe('TorrentControlsComponent', () => {
  let component: TorrentControlsComponent;
  let fixture: ComponentFixture<TorrentControlsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TorrentControlsComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TorrentControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
