// Singleton for managing token refresh to prevent multiple simultaneous refresh calls
class TokenRefreshManager {
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;

  public async refreshToken(): Promise<void> {
    // If already refreshing, wait for that operation to complete
    if (this.isRefreshing && this.refreshPromise) {
      await this.refreshPromise;
      return;
    }

    // Start a new refresh operation
    this.isRefreshing = true;
    this.refreshPromise = this.performRefresh();

    try {
      await this.refreshPromise;
    } finally {
      this.isRefreshing = false;
      this.refreshPromise = null;
    }
  }

  private async performRefresh(): Promise<void> {
    const response = await fetch(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, {
      method: 'POST',
    });

    if (!response.ok) {
      // On refresh failure, redirect to home page
      globalThis.location.reload();
      throw new Error('Token refresh failed');
    }
  }
}

export const tokenRefreshManager = new TokenRefreshManager();
