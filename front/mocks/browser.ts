import { setupWorker } from "msw/browser";
import { getAuthStatus200 } from "./handlers/authStatusHandlers";

export const worker = setupWorker(getAuthStatus200);
