/**
 * Validates if a URL string is a safe internal URL.
 * 
 * This function prevents open redirect vulnerabilities by ensuring:
 * 1. The URL is a relative path starting with a single '/'
 * 2. The URL does not start with '//' (protocol-relative URL)
 * 3. The URL does not contain backslashes (prevents path traversal and normalization bypasses)
 * 
 * @param url - The URL string to validate
 * @returns true if the URL is a safe internal URL, false otherwise
 * 
 * @example
 * isInternalUrl('/dashboard') // true
 * isInternalUrl('/dashboard/profile') // true
 * isInternalUrl('//evil.com') // false
 * isInternalUrl('http://evil.com') // false
 * isInternalUrl('https://evil.com') // false
 * isInternalUrl('javascript:alert(1)') // false
 * isInternalUrl('/\\\\evil.com') // false
 */
export const isInternalUrl = (url: string): boolean => {
  // Check if URL starts with a single '/' but not '//'
  if (!url.startsWith('/') || url.startsWith('//')) {
    return false;
  }

  // Reject URLs containing backslashes to prevent path traversal and URL normalization bypasses
  return !url.includes('\\');
};
