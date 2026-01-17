import { server } from '@Mocks/server';
import { describe, expect, it, vi, beforeEach, afterEach, beforeAll, afterAll } from 'vitest';
import { tokenRefreshManager } from '../tokenRefreshManager';

const mockFetch = vi.fn();

beforeAll(() => {
  server.close();
  global.fetch = mockFetch;
});

afterAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
  vi.restoreAllMocks();
});

describe('TokenRefreshManager', () => {
  let mockLocationReload: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockLocationReload = vi.fn();

    Object.defineProperty(globalThis, 'location', {
      value: {
        reload: mockLocationReload,
      },
      writable: true,
    });

    import.meta.env.FRONT_BACKEND_URL = 'http://localhost:3000/api';
  });

  describe('refreshToken', () => {
    it('should successfully refresh token when API returns ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
      });

      await tokenRefreshManager.refreshToken();

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:3000/api/v1/auth/refresh', {
        method: 'POST',
      });
      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should reload page and throw error when API returns error response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
      });

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Token refresh failed');

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:3000/api/v1/auth/refresh', {
        method: 'POST',
      });
      expect(mockLocationReload).toHaveBeenCalledTimes(1);
    });

    it('should handle multiple simultaneous refresh requests without duplicate API calls', async () => {
      mockFetch.mockImplementation(() => new Promise(resolve => setTimeout(() => resolve({ ok: true }), 100)));

      const refreshPromises = [
        tokenRefreshManager.refreshToken(),
        tokenRefreshManager.refreshToken(),
        tokenRefreshManager.refreshToken(),
      ];

      await Promise.all(refreshPromises);

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should handle subsequent refresh requests after first one completes', async () => {
      mockFetch.mockResolvedValue({ ok: true });

      await tokenRefreshManager.refreshToken();
      await tokenRefreshManager.refreshToken();

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should handle fetch network error', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Network error');

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockLocationReload).not.toHaveBeenCalled();
    });

    it('should reset refresh state even when performRefresh throws error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
      });

      await expect(tokenRefreshManager.refreshToken()).rejects.toThrow('Token refresh failed');

      mockFetch.mockResolvedValueOnce({ ok: true });

      await tokenRefreshManager.refreshToken();

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(mockLocationReload).toHaveBeenCalledTimes(1);
    });
  });
});
