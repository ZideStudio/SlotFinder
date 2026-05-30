import { setupWorker } from "msw/browser";
import { getAuthStatus403 } from "./handlers/authStatusHandlers";

export const worker = setupWorker(getAuthStatus403);
