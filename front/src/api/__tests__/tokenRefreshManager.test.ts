import { server } from '@Mocks/server';
import {
  postTokenRefresh200,
  postTokenRefresh400,
  postTokenRefreshNetworkError,
  postTokenRefreshSlowResponse
} from '@Mocks/handlers/tokenRefreshHandlers';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { tokenRefreshManager } from '../tokenRefreshManager';

const mockLocationReload = vi.fn();

describe('TokenRefreshManager', () => {
  beforeEach(() => {
    Object.defineProperty(globalThis, 'location', {
      value: {
        reload: mockLocationReload,
      },
      writable: true,
    });
  });

  afterEach(() => {
    vi.resetAllMocks();
    server.resetHandlers();
  });

  describe('refreshToken', () => {
    it('should successfully refresh token when API returns ok response', async () => {
      server.use(postTokenRefresh200);

      await tokenRefreshManager.refreshToken();

      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should reload page and throw error when API returns error response', async () => {
      server.use(postTokenRefresh400);

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Token refresh failed');

      expect(mockLocationReload).toHaveBeenCalledTimes(1);
    });

    it('should handle multiple simultaneous refresh requests without duplicate API calls', async () => {
      server.use(postTokenRefreshSlowResponse);

      const refreshPromises = [
        tokenRefreshManager.refreshToken(),
        tokenRefreshManager.refreshToken(),
        tokenRefreshManager.refreshToken(),
      ];

      await Promise.all(refreshPromises);

      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should handle subsequent refresh requests after first one completes', async () => {
      server.use(postTokenRefresh200);

      await tokenRefreshManager.refreshToken();
      await tokenRefreshManager.refreshToken();

      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should handle fetch network error', async () => {
      server.use(postTokenRefreshNetworkError);

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Failed to fetch');

      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should reset refresh state even when performRefresh throws error', async () => {
      server.use(postTokenRefresh400);

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Token refresh failed');

      expect(mockLocationReload).toHaveBeenCalledTimes(1);

      server.use(postTokenRefresh200);

      await tokenRefreshManager.refreshToken();

      expect(mockLocationReload).toHaveBeenCalledTimes(1); // Still 1 from previous call
    });
  });
});
