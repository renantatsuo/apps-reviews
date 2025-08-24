/**
 * Custom error class for HTTP errors.
 */
export class HttpError extends Error {
  /**
   * Constructor.
   * @param status the HTTP status code
   * @param message the error message
   */
  constructor(public status: number, message: string) {
    super(message);
    this.status = status;
  }
}

/**
 * Checks if the error is an HttpError.
 * @param error the error to check
 * @returns true if the error is an HttpError
 */
export const isHttpError = (error: unknown): error is HttpError => {
  return error instanceof HttpError;
};
