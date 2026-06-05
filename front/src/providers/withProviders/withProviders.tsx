import { QueryClientProvider } from "@Front/providers/QueryClientProvider";
import { type ComponentProps, type ComponentType, createElement } from "react";
import { AuthenticationContextProvider } from "../../contexts/AuthenticationContext/AuthenticationContextProvider";
import { ToastProvider } from "../../ui/utils/toast/toastProvider/ToastProvider";
import { LoaderProvider } from "../loaderProvider/LoaderProvider";

export const withProvider = <WithProviderProps extends object>(
  Component: ComponentType<WithProviderProps>,
  queryClient: ComponentProps<typeof QueryClientProvider>["client"],
) => {
  const WithProvider = (props: WithProviderProps) =>
    createElement(
      QueryClientProvider,
      { client: queryClient },
      <AuthenticationContextProvider>
        <LoaderProvider>
          <ToastProvider>
            <Component {...props} />
          </ToastProvider>
        </LoaderProvider>
      </AuthenticationContextProvider>,
    );

  WithProvider.displayName = `withProvider(${Component.displayName || Component.name})`;

  return WithProvider;
};
