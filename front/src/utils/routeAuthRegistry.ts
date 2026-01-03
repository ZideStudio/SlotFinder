/**
 * Registry for tracking which routes require authentication
 * Routes register themselves by calling registerUnauthenticatedPath() in their route definition files
 */
class RouteAuthRegistry {
  private unauthenticatedPaths: Set<string> = new Set();

  /**
   * Register a path that does not require authentication
   * Should be called in route definition files for routes with handle.mustBeAuthenticate = false
   */
  registerUnauthenticatedPath(path: string): void {
    this.unauthenticatedPaths.add(path);
  }

  /**
   * Check if a path requires authentication
   * @param pathname - The current pathname
   * @returns true if authentication is required, false otherwise
   */
  requiresAuthentication(pathname: string): boolean {
    // Check exact matches first
    if (this.unauthenticatedPaths.has(pathname)) {
      return false;
    }

    // Check if any registered unauthenticated path is a prefix of the current path
    for (const path of this.unauthenticatedPaths) {
      if (pathname.startsWith(`/${path}`) || pathname.startsWith(path)) {
        return false;
      }
    }

    // Default: require authentication
    return true;
  }

  /**
   * Get all registered unauthenticated paths (for testing/debugging)
   */
  getUnauthenticatedPaths(): string[] {
    return Array.from(this.unauthenticatedPaths);
  }
}

// Export singleton instance
export const routeAuthRegistry = new RouteAuthRegistry();
