import { FileSizeFormatterPipe } from './file-size-formatter.pipe';

describe('FileSizeFormatterPipe', () => {
  it('create an instance', () => {
    const pipe = new FileSizeFormatterPipe();
    expect(pipe).toBeTruthy();
  });
});
