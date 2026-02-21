import { QueryClientProvider } from '@Front/providers/QueryClientProvider';
import type { ComponentProps, ComponentType } from 'react';
import { AuthenticationContextProvider } from '../../contexts/AuthenticationContext/AuthenticationContextProvider';
import { ToastProvider } from '../ToastProvider/ToastProvider';

type WithRootProps = {
  queryClient: ComponentProps<typeof QueryClientProvider>['client'];
};

export const withProvider = <WithProviderProps extends object>(Component: ComponentType<WithProviderProps>) => {
  const WithProvider = ({ queryClient, ...props }: WithProviderProps & WithRootProps) => (
    <QueryClientProvider client={queryClient}>
      <AuthenticationContextProvider>
        <ToastProvider>
          <Component {...(props as WithProviderProps)} />
        </ToastProvider>
      </AuthenticationContextProvider>
    </QueryClientProvider>
  );

  WithProvider.displayName = `withProvider(${Component.displayName || Component.name})`;

  return WithProvider;
};
