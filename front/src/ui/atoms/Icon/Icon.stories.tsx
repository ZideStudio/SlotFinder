import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import GoogleIcon from '@Front/assets/svg/google_icon.svg?react';
import { Icon } from './Icon';

const meta = {
  component: Icon,
} satisfies Meta<typeof Icon>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    icon: GoogleIcon,
  },
};
