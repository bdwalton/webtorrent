import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'fileSizeFormatter'
})
export class FileSizeFormatterPipe implements PipeTransform {

  transform(value: unknown, ...args: unknown[]): unknown {
    return null;
  }

}
