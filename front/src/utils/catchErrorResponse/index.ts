export const getFormattedError = (error: Error): Error => {
  try {
    const parsedError = JSON.parse(error.message);
    if (parsedError && typeof parsedError === 'object' && 'message' in parsedError) {
      return new Error(parsedError.message);
    }
  } catch {}

  return new Error('An unexpected error occurred');
};
