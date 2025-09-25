export const getFormattedErrorMessage = (error: Error | null): string | undefined => {
  if (!error) {
    return undefined;
  }

  try {
    const parsedError = JSON.parse(error.message);
    if (parsedError && typeof parsedError === 'object' && 'message' in parsedError) {
      return parsedError.message;
    }
  } catch {
    // If parsing fails, return the fallback message
  }

  return 'An unexpected error occurred';
};
