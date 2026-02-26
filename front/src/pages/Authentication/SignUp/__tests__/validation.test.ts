import { describe, expect, it } from 'vitest';
import { EMAIL_REGEX } from '../constants';

describe('SignUp validation', () => {
  describe('EMAIL_REGEX', () => {
    it.each([
      'user+alias@example.com',
      'test+tag@domain.org',
      'name+test+multiple@test.com',
      'simple@example.com',
      'user.name@example.com',
      'user_name@example.com',
      'user-name@example.com',
      'user123@example.com',
      'user@sub.domain.com',
      'u@example.co',
      'user%test@example.com',
      'user%test_123@example.com',
    ])('should accept valid email: %s', email => {
      expect(EMAIL_REGEX.test(email)).toBe(true);
    });

    it.each([
      '',
      'notanemail',
      '@example.com',
      'user@',
      'user@domain',
      'user@@example.com',
      'user @example.com',
      'user@exam ple.com',
      'user@.com',
      'user@domain.',
    ])('should reject invalid email: %s', email => {
      expect(EMAIL_REGEX.test(email)).toBe(false);
    });
  });
});
