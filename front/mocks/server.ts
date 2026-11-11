import { setupServer } from "msw/node";
import { defaultHandlers } from "./handlers/defaultHandlers";

export const server = setupServer(...defaultHandlers);
