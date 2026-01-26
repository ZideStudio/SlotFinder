import { AuthenticationContextProvider } from '@Front/contexts/AuthenticationContext/AuthenticationContextProvider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook } from '@testing-library/react';
import { type PropsWithChildren } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { useAuthenticationContext } from '../useAuthenticationContext';

describe('useAuthenticationContext', () => {
  it('should return the context value when used within AuthenticationContextProvider', () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
          refetchOnWindowFocus: false,
          refetchOnReconnect: false,
        },
      },
    });

    const wrapper = ({ children }: PropsWithChildren) => (
      <QueryClientProvider client={queryClient}>
        <AuthenticationContextProvider>{children}</AuthenticationContextProvider>
      </QueryClientProvider>
    );

    const { result } = renderHook(() => useAuthenticationContext(), {
      wrapper,
    });

    expect(result.current).toBeDefined();
    expect(result.current).toHaveProperty('isAuthenticated');
    expect(result.current).toHaveProperty('checkAuthentication');
    expect(result.current).toHaveProperty('postAuthRedirectPath');
    expect(result.current).toHaveProperty('setPostAuthRedirectPath');
    expect(result.current).toHaveProperty('resetPostAuthRedirectPath');
  });

  it('should throw an error when used outside of AuthenticationContextProvider', () => {
    // Suppress console.error for this test since we expect an error
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      renderHook(() => useAuthenticationContext());
    }).toThrow('useAuthenticationContext must be used within a AuthenticationContextProvider');

    consoleErrorSpy.mockRestore();
  });
});
