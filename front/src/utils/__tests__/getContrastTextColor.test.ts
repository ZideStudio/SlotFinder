import { getContrastTextColor } from '../getContrastTextColor';

describe('GetContrastTextColor', () => {
  it('should return #FFFFFF', () => {
    const res = getContrastTextColor('#000000');
    expect(res).toBe('#FFFFFF');
  });

  it('should return #000000', () => {
    const res = getContrastTextColor('#FFFFFF');
    expect(res).toBe('#000000');
  });
});
