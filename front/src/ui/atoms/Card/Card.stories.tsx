import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { Card } from './Card';

const meta = {
  title: 'Atoms/Card',
  component: Card,
  args: {
    children: 'Card',
    className: 'custom-class',
  },
  argTypes: {
    children: {
      control: { type: 'text' },
    },
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Card>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: 'Card',
  },
};
