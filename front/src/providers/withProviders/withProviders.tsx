import { QueryClientProvider } from '@Front/providers/QueryClientProvider';
import type { ComponentProps, ComponentType } from 'react';
import { AuthenticationContextProvider } from '../../contexts/AuthenticationContext/AuthenticationContextProvider';

type WithRootProps = {
  queryClient: ComponentProps<typeof QueryClientProvider>['client'];
};

export const withProvider = <ComponentProps extends object>(Component: ComponentType<ComponentProps>) => {
  const WithProvider = ({ queryClient, ...props }: ComponentProps & WithRootProps) => (
    <QueryClientProvider client={queryClient}>
      <AuthenticationContextProvider>
        <Component {...(props as ComponentProps)} />
      </AuthenticationContextProvider>
    </QueryClientProvider>
  );

  WithProvider.displayName = `withProvider(${Component.displayName || Component.name})`;

  return WithProvider;
};
