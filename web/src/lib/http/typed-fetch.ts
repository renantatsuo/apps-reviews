import { HttpError } from "~/lib/http/http-error";

/**
 * A "Typed" fetch implementation of Http.
 */

export const execute = async <T>(req: Request): Promise<T> => {
  const response = await fetch(req);
  if (!response.ok) {
    throw new HttpError(response.status, response.statusText);
  }
  try {
    const asJson = await response.json();
    return asJson as T;
  } catch (_) {
    throw new HttpError(500, "Failed to parse JSON response");
  }
};

export const get = async <T>(url: string, opts?: RequestInit): Promise<T> => {
  return execute<T>(new Request(url, { ...opts, method: "GET" }));
};
