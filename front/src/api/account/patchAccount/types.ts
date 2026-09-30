import type { ErrorResponseCodeType } from "@Front/types/api.types";

export type PatchAccountErrorCodeType = ErrorResponseCodeType<
  "USERNAME_ALREADY_TAKEN" | "INVALID_PASSWORD_FORMAT" | "INVALID_COLOR_FORMAT"
>;
