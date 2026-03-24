import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { Heading } from './Heading';

const meta = {
  component: Heading,
} satisfies Meta<typeof Heading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    level: 1,
    children: 'Heading 1',
    className: 'custom-class',
  },
};

export const Level2: Story = {
  args: {
    level: 2,
    children: 'Heading 2',
    className: 'custom-class',
  },
};

export const Level3: Story = {
  args: {
    level: 3,
    children: 'Heading 3',
    className: 'custom-class',
  },
};
