import { setupWorker } from "msw/browser";
import {
  getAccountAvatar200,
  getAccountMe200,
} from "./handlers/accountHandlers";
import { getAuthStatus403 } from "./handlers/authStatusHandlers";

export const worker = setupWorker(
  getAuthStatus403("USERNAME_MISSING"),
  getAccountMe200,
  getAccountAvatar200,
);
