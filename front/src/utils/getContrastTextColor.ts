// oxlint-disable no-magic-numbers
export const getContrastTextColor = (hex: string): '#000000' | '#FFFFFF' => {
  const cleanHex = hex.replace('#', '');

  const r = Number.parseInt(cleanHex.substring(0, 2), 16) / 255;
  const g = Number.parseInt(cleanHex.substring(2, 4), 16) / 255;
  const b = Number.parseInt(cleanHex.substring(4, 6), 16) / 255;

  const luminance = 0.2126 * r + 0.7152 * g + 0.0722 * b;

  return luminance > 0.5 ? '#000000' : '#FFFFFF';
};
