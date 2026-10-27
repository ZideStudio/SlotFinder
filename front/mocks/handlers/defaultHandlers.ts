import { getAccountAvatar200, getAccountMe200 } from "./accountHandlers";
import { getAuthStatus200 } from "./authStatusHandlers";

// Signed-in baseline shared by Node, browser and Playwright; override per test.
export const defaultHandlers = [
  getAuthStatus200,
  getAccountMe200,
  getAccountAvatar200,
];
