import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { Link } from './Link';

const meta = {
  title: 'Atoms/Link',
  component: Link,
  args: {
    href: 'https://react.dev',
    children: 'Lien par défaut',
    className: 'custom-class',
  },
  argTypes: {
    openInNewTab: {
      control: { type: 'boolean' },
    },
    onClick: {
      action: true,
      table: { disable: true },
    },
  },
} satisfies Meta<typeof Link>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const OpenInNewTab: StoryObj<typeof meta> = {
  args: {
    openInNewTab: true,
    children: 'Lien ouvert dans un nouvel onglet',
  },
};
