import { useToastContext } from '../useToastContext';

describe('ToastContext', () => {
  it('should throw an error if used outside of ToastProvider', () => {
    const consoleErrorSpy = vitest.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => useToastContext()).toThrow("Cannot read properties of null (reading 'useContext')");

    consoleErrorSpy.mockRestore();
  });
});
