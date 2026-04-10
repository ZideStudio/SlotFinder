import SendIcon from '@material-symbols/svg-400/outlined/send.svg?react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { ClickIcon } from './ClickIcon';

const meta = {
  title: 'Molecules/ClickIcon',
  component: ClickIcon,
} satisfies Meta<typeof ClickIcon>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    icon: SendIcon,
    onClick: { action: true, table: { disable: true } },
  },
};
