import { describe, expect, it } from 'vitest';
import { isInternalUrl } from '../isInternalUrl';

describe('isInternalUrl', () => {
  describe('should return true for valid internal URLs', () => {
    it('should accept various valid internal paths', () => {
      expect(isInternalUrl('/dashboard')).toBe(true);
      expect(isInternalUrl('/dashboard/profile')).toBe(true);
      expect(isInternalUrl('/dashboard?tab=settings')).toBe(true);
      expect(isInternalUrl('/dashboard#section')).toBe(true);
      expect(isInternalUrl('/dashboard?tab=settings#section')).toBe(true);
      expect(isInternalUrl('/')).toBe(true);
    });
  });

  describe('should return false for invalid or external URLs', () => {
    it('should reject protocol-relative URLs', () => {
      expect(isInternalUrl('//evil.com')).toBe(false);
      expect(isInternalUrl('//evil.com/dashboard')).toBe(false);
    });

    it('should reject absolute URLs with protocols', () => {
      expect(isInternalUrl('http://evil.com')).toBe(false);
      expect(isInternalUrl('https://evil.com')).toBe(false);
      expect(isInternalUrl('javascript:alert(1)')).toBe(false);
      expect(isInternalUrl('data:text/html,<script>alert(1)</script>')).toBe(false);
    });

    it('should reject URLs with backslashes', () => {
      expect(isInternalUrl('\\evil.com')).toBe(false);
      expect(isInternalUrl('/\\evil.com')).toBe(false);
      expect(isInternalUrl('/\\\\evil.com')).toBe(false);
    });

    it('should reject invalid path formats', () => {
      expect(isInternalUrl('dashboard')).toBe(false);
      expect(isInternalUrl('')).toBe(false);
    });
  });
});
