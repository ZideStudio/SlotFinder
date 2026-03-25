import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import SendIcon from '@material-symbols/svg-400/outlined/send.svg?react';
import { Button } from './Button';

const meta = {
  title: 'Molecules/Inputs/Button',
  component: Button,
  args: {
    children: 'Send message',
    variant: 'primary',
    color: 'default',
    disabled: false,
  },
  argTypes: {
    children: { control: 'text' },
    variant: { control: 'select', options: ['primary', 'secondary'] },
    color: { control: 'select', options: ['default', 'neutral', 'danger'] },
    disabled: { control: 'boolean' },
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
  render: args => (
    <Button {...args} icon={SendIcon}>
      {args.children}
    </Button>
  ),
};
