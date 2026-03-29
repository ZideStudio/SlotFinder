// oxlint-disable no-magic-numbers
import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { Heading } from './Heading';

const meta = {
  component: Heading,
  title: 'Atoms/Heading',
  argTypes: {
    level: {
      control: { type: 'select', options: [1, 2, 3] },
    },
  },
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
