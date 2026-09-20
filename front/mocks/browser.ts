import { setupWorker } from "msw/browser";
import {
  getAccountAvatar200,
  getAccountMe200,
} from "./handlers/accountHandlers";
import { getAuthStatus200 } from "./handlers/authStatusHandlers";

export const worker = setupWorker(
  getAuthStatus200,
  getAccountMe200,
  getAccountAvatar200,
);
