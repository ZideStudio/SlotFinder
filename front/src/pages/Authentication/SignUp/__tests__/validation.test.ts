import { describe, expect, it } from 'vitest';
import { EMAIL_REGEX } from '../constants';

describe('SignUp validation', () => {
  describe('EMAIL_REGEX', () => {
    it('should accept valid email with plus sign', () => {
      const validEmails = [
        'jules+test@zide.fr',
        'user+alias@example.com',
        'test+123@domain.org',
        'name+tag+multiple@test.com',
      ];

      validEmails.forEach(email => {
        expect(EMAIL_REGEX.test(email)).toBe(true);
      });
    });

    it('should accept valid standard email addresses', () => {
      const validEmails = [
        'simple@example.com',
        'user.name@example.com',
        'user_name@example.com',
        'user-name@example.com',
        'user123@example.com',
        'user@sub.domain.com',
        'u@example.co',
      ];

      validEmails.forEach(email => {
        expect(EMAIL_REGEX.test(email)).toBe(true);
      });
    });

    it('should accept email with percent and underscore characters', () => {
      const validEmails = ['user%test@example.com', 'user_test@example.com', 'user%test_123@example.com'];

      validEmails.forEach(email => {
        expect(EMAIL_REGEX.test(email)).toBe(true);
      });
    });

    it('should reject invalid email addresses', () => {
      const invalidEmails = [
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
      ];

      invalidEmails.forEach(email => {
        expect(EMAIL_REGEX.test(email)).toBe(false);
      });
    });
  });
});
