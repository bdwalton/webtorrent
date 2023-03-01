import { Pipe, PipeTransform } from '@angular/core';

const FILE_SIZE_UNITS = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

@Pipe({
  name: 'fileSizeFormatter'
})
export class FileSizeFormatterPipe implements PipeTransform {
  transform(sizeInBytes: number): string {

    if (sizeInBytes == 0) {
      return `0 B`
    }

    let power = Math.round(Math.log(sizeInBytes) / Math.log(1024));
    power = Math.min(power, FILE_SIZE_UNITS.length - 1);

    const size = sizeInBytes / Math.pow(1024, power); // size in new units
    const formattedSize = Math.round(size * 100) / 100; // keep up to 2 decimals
    const unit = FILE_SIZE_UNITS[power];

    return `${formattedSize} ${unit}`;
  }
}
