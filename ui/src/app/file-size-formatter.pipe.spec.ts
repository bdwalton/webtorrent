import { FileSizeFormatterPipe } from './file-size-formatter.pipe';

describe('FileSizeFormatterPipe', () => {
  it('transforms 518285851 to "0.48 GiB"', () => {
    const pipe = new FileSizeFormatterPipe();
    expect(pipe.transform(518285851)).toBe("0.48 GiB");
  });

  it('transforms 0 to "0 B"', () => {
    const pipe = new FileSizeFormatterPipe();
    expect(pipe.transform(51)).toBe("51 B");
  });

  it('transforms 1 to "1 B"', () => {
    const pipe = new FileSizeFormatterPipe();
    expect(pipe.transform(51)).toBe("51 B");
  });

  it('transforms 1024 to "1 KiB"', () => {
    const pipe = new FileSizeFormatterPipe();
    expect(pipe.transform(1024)).toBe("1024 KiB");
  });
});
