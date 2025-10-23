import DiscordIcon from '@Front/assets/svg/discord_icon.svg?react';
import GithubIcon from '@Front/assets/svg/github_icon.svg?react';
import GoogleIcon from '@Front/assets/svg/google_icon.svg?react';
import type { OAuthProvider } from './types';

export const oauthProvidersData: Omit<OAuthProvider, 'href'>[] = [
  {
    id: 'google',
    label: 'Google',
    ariaLabel: 'signInWithGoogle',
    icon: <GoogleIcon width={24} height={24} aria-hidden />,
  },
  {
    id: 'github',
    label: 'GitHub',
    ariaLabel: 'signInWithGithub',
    icon: <GithubIcon width={24} height={24} aria-hidden />,
  },
  {
    id: 'discord',
    label: 'Discord',
    ariaLabel: 'signInWithDiscord',
    icon: <DiscordIcon width={24} height={24} aria-hidden />,
  },
];
