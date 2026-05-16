import { renderHook } from '@testing-library/react';
import { usePopover } from '../usePopover';

describe('usePopover', () => {
  describe('triggerProps', () => {
    it('should set popoverTarget matching popoverProps.id', () => {
      const { result } = renderHook(() => usePopover());

      expect(result.current.triggerProps.popoverTarget).toBe(result.current.popoverProps.id);
    });

    it('should set data-popover-trigger attribute', () => {
      const { result } = renderHook(() => usePopover());

      expect(result.current.triggerProps['data-popover-trigger']).toBe('');
    });

    it('should set the anchor name CSS custom property in style', () => {
      const { result } = renderHook(() => usePopover());
      const anchorName = (result.current.triggerProps.style as Record<string, string>)['--popover-anchor-name'];

      expect(anchorName).toMatch(/^--popover-/u);
      expect(anchorName).not.toContain(':');
    });
  });

  describe('popoverProps', () => {
    it('should set id to a non-empty string', () => {
      const { result } = renderHook(() => usePopover());

      expect(result.current.popoverProps.id).toMatch(/^popover-/u);
    });

    it('should share the same anchor name CSS custom property as triggerProps', () => {
      const { result } = renderHook(() => usePopover());
      const triggerAnchor = (result.current.triggerProps.style as Record<string, string>)['--popover-anchor-name'];
      const popoverAnchor = (result.current.popoverProps.style as Record<string, string>)['--popover-anchor-name'];

      expect(triggerAnchor).toBe(popoverAnchor);
    });
  });

  describe('anchor name uniqueness', () => {
    it('should generate a unique anchor name for each hook instance', () => {
      const { result: result1 } = renderHook(() => usePopover());
      const { result: result2 } = renderHook(() => usePopover());

      const anchor1 = (result1.current.triggerProps.style as Record<string, string>)['--popover-anchor-name'];
      const anchor2 = (result2.current.triggerProps.style as Record<string, string>)['--popover-anchor-name'];

      expect(anchor1).not.toBe(anchor2);
    });
  });
});
