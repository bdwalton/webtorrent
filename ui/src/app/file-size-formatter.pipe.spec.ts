import { FileSizeFormatterPipe } from './file-size-formatter.pipe';

describe('FileSizeFormatterPipe', () => {
  const pipe = new FileSizeFormatterPipe();
  it('transforms 0 to "0 B"', () => {
    expect(pipe.transform(0)).toBe("0 B");
  });

  it('transforms 1 to "1 B"', () => {
    expect(pipe.transform(1)).toBe("1 B");
  });

  it('transforms 1024 to "1 KiB"', () => {
    expect(pipe.transform(1024)).toBe("1 KiB");
  });

  it('transforms 518285851 to "0.48 GiB"', () => {
    expect(pipe.transform(518285851)).toBe("0.48 GiB");
  });

  it('transforms 1048576 to "1 MiB"', () => {
    expect(pipe.transform(1048576)).toBe("1 MiB");
  });

  it('transforms 1125899906842624 to "1 PiB"', () => {
    expect(pipe.transform(1125899906842624)).toBe("1 PiB");
  });

  it('transforms 1152921504606846976 to "1 EiB"', () => {
    expect(pipe.transform(1152921504606846976)).toBe("1 EiB");
  });

  it('transforms 1180591620717411303424 to "1 ZiB"', () => {
    expect(pipe.transform(1180591620717411303424)).toBe("1 ZiB");
  });

  it('transforms 1208925819614629174706176 to "1 YiB"', () => {
    expect(pipe.transform(1208925819614629174706176)).toBe("1 YiB");
  });
});
