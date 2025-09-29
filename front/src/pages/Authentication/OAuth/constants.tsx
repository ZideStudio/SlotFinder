import DiscordIcon from '@Front/assets/svg/discord_icon.svg?react';
import GithubIcon from '@Front/assets/svg/github_icon.svg?react';
import GoogleIcon from '@Front/assets/svg/google_icon.svg?react';
import type { ReactNode } from 'react';

type oauthProvider = {
  id: string;
  label: string;
  href: string;
  ariaLabel: 'signInWithGoogle' | 'signInWithGitHub' | 'signInWithDiscord';
  icon: ReactNode;
};

export const oauthProviders: oauthProvider[] = [
  {
    id: 'google',
    label: 'Google',
    href: '#google-oauth',
    ariaLabel: 'signInWithGoogle',
    icon: <GoogleIcon width={24} height={24} aria-hidden />,
  },
  {
    id: 'github',
    label: 'GitHub',
    href: '#github-oauth',
    ariaLabel: 'signInWithGitHub',
    icon: <GithubIcon width={24} height={24} aria-hidden />,
  },
  {
    id: 'discord',
    label: 'Discord',
    href: '#discord-oauth',
    ariaLabel: 'signInWithDiscord',
    icon: <DiscordIcon width={24} height={24} aria-hidden />,
  },
];
