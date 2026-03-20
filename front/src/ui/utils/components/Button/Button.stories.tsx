import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import GoogleIcon from '@Front/assets/svg/google_icon.svg?react';
import { Button } from './Button';

const meta = {
  title: 'Utils/Components/Button',
  component: Button,
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
    <Button {...args} icon={GoogleIcon}>
      Se connecter avec Google
    </Button>
  ),
};
