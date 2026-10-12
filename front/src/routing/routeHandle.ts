export type RouteHandle = {
  mustBeAuthenticate?: boolean;
  hideHeader?: boolean;
};
export type AuthenticationRouteHandle = RouteHandle & { title?: string; }
