/**
 * Checks if the current page requires authentication by checking the URL path
 * This relies on the convention that authenticated pages don't include '/sign-up' or '/oauth'
 * @returns true if the current page likely requires authentication, false otherwise
 */
export const isAuthenticatedPage = (): boolean => {
  const currentPath = window.location.pathname;
  
  // Pages that don't require authentication
  const unauthenticatedPaths = ['/sign-up', '/oauth'];
  
  // Check if the current path matches any unauthenticated path
  const isUnauthenticated = unauthenticatedPaths.some(path => currentPath.includes(path));
  
  // If not on an unauthenticated path, assume authentication is required
  return !isUnauthenticated;
};
