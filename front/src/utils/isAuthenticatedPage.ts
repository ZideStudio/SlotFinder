import { routeAuthRegistry } from './routeAuthRegistry';

/**
 * Checks if the current page requires authentication
 * Uses the route authentication registry which mirrors the route configuration
 * without causing circular dependencies.
 * 
 * This follows the same convention as AuthenticationProtection component:
 * - Pages explicitly marked as unauthenticated (handle.mustBeAuthenticate = false) don't require auth
 * - All other pages require authentication by default
 * 
 * @returns true if the current page requires authentication, false otherwise
 */
export const isAuthenticatedPage = (): boolean => {
  const currentPath = window.location.pathname;
  
  // Remove leading slash for consistent comparison
  const normalizedPath = currentPath.startsWith('/') ? currentPath.substring(1) : currentPath;
  
  return routeAuthRegistry.requiresAuthentication(normalizedPath);
};
