import { setupWorker } from "msw/browser";
import { defaultHandlers } from "./handlers/defaultHandlers";

export const worker = setupWorker(...defaultHandlers);
