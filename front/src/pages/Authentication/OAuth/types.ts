import type { FC, SVGProps } from "react";

export type OAuthProviderName = "discord" | "google" | "github";

export type OAuthProvider = {
  id: OAuthProviderName;
  label: string;
  href: string;
  ariaLabel: `signInWith${Capitalize<OAuthProviderName>}`;
  icon: FC<SVGProps<SVGSVGElement>>;
};
  