export type RouteHandle = {
  mustBeAuthenticate?: boolean;
  hideHeader?: boolean;
  hideOAuth?: boolean;
};
export type AuthenticationRouteHandle = RouteHandle & { title?: string };
