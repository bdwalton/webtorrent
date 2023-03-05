import { Component, Input } from '@angular/core';
import { Progress } from '../torrent.service';

@Component({
  selector: 'app-progress-component',
  templateUrl: './progress-component.component.html',
  styleUrls: ['./progress-component.component.scss'],
})
export class ProgressComponentComponent {
  @Input() prg = new Progress(); // A torrent or file progress
}
