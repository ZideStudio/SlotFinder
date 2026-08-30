export const bytesToMB = (sizeInBytes: number): number =>
  Number.parseFloat((sizeInBytes / (1024 * 1024)).toFixed(2));

export const mbToBytes = (sizeInMB: number): number => sizeInMB * 1024 * 1024;
