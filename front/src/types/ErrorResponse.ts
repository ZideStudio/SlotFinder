export class ErrorResponse<ErrorCodeType extends string> extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ErrorResponse';
  }

  getErrorCode(): ErrorCodeType | 'SERVER_ERROR' {
    try {
      const parsed = JSON.parse(this.message);
      if (parsed && typeof parsed === 'object' && 'code' in parsed && typeof parsed.code === 'string') {
        return parsed.code;
      }
    } catch {
      // ignore
    }
    return 'SERVER_ERROR';
  }
}
