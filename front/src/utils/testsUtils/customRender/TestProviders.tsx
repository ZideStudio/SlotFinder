import { AuthenticationContextProvider } from "@Front/contexts/AuthenticationContext/AuthenticationContextProvider";
import { LoaderProvider } from "@Front/providers/loaderProvider/LoaderProvider";
import { ToastProvider } from "@Front/ui/utils/toast/toastProvider/ToastProvider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { PropsWithChildren } from "react";

export const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
      },
    },
  });

type TestProvidersProps = PropsWithChildren<{
  client?: QueryClient;
}>;

export const TestProviders = ({
  children,
  client = createQueryClient(),
}: TestProvidersProps) => (
  <QueryClientProvider client={client}>
    <LoaderProvider>
      <ToastProvider>
        <AuthenticationContextProvider>
          {children}
        </AuthenticationContextProvider>
      </ToastProvider>
    </LoaderProvider>
  </QueryClientProvider>
);
