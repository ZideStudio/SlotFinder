import { describe, expect, it } from 'vitest';
import { getFormattedError } from '..';

describe('getFormattedError', () => {
  it('returns parsed message from valid JSON error', () => {
    const error = new Error(JSON.stringify({ message: 'Custom error' }));
    const result = getFormattedError(error);
    expect(result).toBeInstanceOf(Error);
    expect(result.message).toBe('Custom error');
  });

  it('returns default message for invalid JSON error', () => {
    const error = new Error('not a json');
    const result = getFormattedError(error);
    expect(result).toBeInstanceOf(Error);
    expect(result.message).toBe('An unexpected error occurred');
  });

  it('returns fallback message for JSON without message property', () => {
    const error = new Error(JSON.stringify({ foo: 'bar' }));
    const result = getFormattedError(error);
    expect(result).toBeInstanceOf(Error);
    expect(result.message).toBe('An unexpected error occurred');
  });

  it('returns fallback message if error.message is empty', () => {
    const error = new Error('');
    const result = getFormattedError(error);
    expect(result).toBeInstanceOf(Error);
    expect(result.message).toBe('An unexpected error occurred');
  });
});
