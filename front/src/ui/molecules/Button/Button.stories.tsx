import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import PaletteIcon from '@material-symbols/svg-400/outlined/palette.svg?react';
import { Button } from './Button';

const meta = {
  title: 'Molecules/Button',
  component: Button,
  args: {
    children: 'Send message',
    variant: 'primary',
    color: 'default',
    disabled: false,
    isLoading: false,
  },
  argTypes: {
    children: { control: 'text' },
    variant: { control: 'select', options: ['primary', 'secondary'] },
    color: { control: 'select', options: ['default', 'neutral', 'danger'] },
    disabled: { control: 'boolean' },
    isLoading: { control: 'boolean' },
    as: { table: { disable: true } },
    icon: { table: { disable: true } },
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    icon: PaletteIcon,
  },
};

export const Loading: Story = {
  args: {
    isLoading: true,
  },
};
