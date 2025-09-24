export type JsonValue = JsonObject | JsonArray | string | number | boolean | null;
export type JsonArray = JsonValue[];
export type JsonObject = {
  // oxlint-disable-next-line consistent-indexed-object-style
  [key: string]: JsonValue;
};
export type Json = JsonObject | JsonArray;
