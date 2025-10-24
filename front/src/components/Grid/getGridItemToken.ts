import type { Breakpoint, ColSpan, ColStart } from './types';

const breakpoints: Breakpoint[] = ['mobile', 'tablet', 'desktop-small', 'desktop-medium', 'desktop-large'];

export const getGridItemToken = (
  prefix: 'col' | 'start',
  data: ColSpan | ColStart | Partial<Record<Breakpoint, ColSpan | ColStart>>,
) => {
  const tokens: Record<string, ColSpan | ColStart> = {};
  if (typeof data === 'number' || typeof data === 'string') {
    tokens[`--${prefix}`] = data;
  } else {
    const sortedBreakpointsConfiguration = Object.entries(data).sort(
      ([breakpointA], [breakpointB]) =>
        breakpoints.indexOf(breakpointA as Breakpoint) - breakpoints.indexOf(breakpointB as Breakpoint),
    );

    for (const [breakpointKey, breakpointValue] of sortedBreakpointsConfiguration) {
      if (breakpointValue) {
        tokens[`--${prefix}`] = breakpointValue;
        tokens[`--${prefix}-${breakpointKey}`] = breakpointValue;
      }
    }
  }

  return tokens;
};
