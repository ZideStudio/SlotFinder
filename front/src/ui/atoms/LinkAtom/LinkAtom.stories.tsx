import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { LinkAtom } from './LinkAtom';

const meta = {
  title: 'Atoms/Links/LinkAtom',
  component: LinkAtom,
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
  decorators: [
    Story => (
      <div style={{ padding: '16px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof LinkAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const OpenInNewTab: StoryObj<typeof meta> = {
  args: {
    openInNewTab: true,
    children: 'Lien ouvert dans un nouvel onglet',
  },
};
