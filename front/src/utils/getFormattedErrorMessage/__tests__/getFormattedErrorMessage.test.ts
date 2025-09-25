import { describe, expect, it } from 'vitest';
import { getFormattedErrorMessage } from '../index';

describe('getFormattedErrorMessage', () => {
  it('should return undefined if error is null', () => {
    expect(getFormattedErrorMessage(null)).toBeUndefined();
  });

  it('should return the message if error.message is a valid JSON with a message property', () => {
    const error = new Error(JSON.stringify({ message: 'Custom error message' }));
    expect(getFormattedErrorMessage(error)).toBe('Custom error message');
  });

  it('should return the fallback message if error.message is not a valid JSON', () => {
    const error = new Error('Not a JSON');
    expect(getFormattedErrorMessage(error)).toBe('An unexpected error occurred');
  });

  it('should return the fallback message if error.message is a JSON without message property', () => {
    const error = new Error(JSON.stringify({ foo: 'bar' }));
    expect(getFormattedErrorMessage(error)).toBe('An unexpected error occurred');
  });

  it('should return the fallback message if error is an empty Error', () => {
    const error = new Error();
    expect(getFormattedErrorMessage(error)).toBe('An unexpected error occurred');
  });
});
