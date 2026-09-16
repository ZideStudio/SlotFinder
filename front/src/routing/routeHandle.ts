export type RouteHandle = {
  mustBeAuthenticate?: boolean;
  hideHeader?: boolean;
  title?: string;
};
export type AuthenticationRouteHandle = RouteHandle & { title?: string };
