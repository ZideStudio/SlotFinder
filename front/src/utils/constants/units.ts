export type DurationUnit = "days" | "hours" | "minutes";

export const UNITS: DurationUnit[] = ["days", "hours", "minutes"];

export const UNIT_LIMITS: Record<DurationUnit, { min: number; max: number }> = {
  days: { min: 0, max: 21 },
  hours: { min: 0, max: 23 },
  minutes: { min: 0, max: 59 },
};
