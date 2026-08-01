import DiscordIcon from "@Front/assets/svg/discord_icon.svg";
import GithubIcon from "@Front/assets/svg/github_icon.svg";
import GoogleIcon from "@Front/assets/svg/google_icon.svg";
import type { OAuthProvider } from "./types";

export const oauthProvidersData: Omit<OAuthProvider, "href">[] = [
  {
    id: "google",
    label: "Google",
    ariaLabel: "signInWithGoogle",
    icon: GoogleIcon,
  },
  {
    id: "github",
    label: "GitHub",
    ariaLabel: "signInWithGithub",
    icon: GithubIcon,
  },
  {
    id: "discord",
    label: "Discord",
    ariaLabel: "signInWithDiscord",
    icon: DiscordIcon,
  },
];
