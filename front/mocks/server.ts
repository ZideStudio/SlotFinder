import { setupServer } from "msw/node";
import { getAccountAvatar200, getAccountMe200 } from "./handlers/accountHandlers";
import { getAuthStatus200 } from "./handlers/authStatusHandlers";

export const server = setupServer(getAuthStatus200, getAccountMe200, getAccountAvatar200);
