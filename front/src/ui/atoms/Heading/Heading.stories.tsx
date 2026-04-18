import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { Heading } from './Heading';

const meta = {
  component: Heading,
  title: 'Atoms/Heading',
  args: {
    level: 1,
  },
  argTypes: {
    level: {
      control: {
        type: 'select',
        // oxlint-disable-next-line no-magic-numbers
        options: [1, 2, 3],
      },
    },
  },
} satisfies Meta<typeof Heading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: 'Heading 1',
    className: 'custom-class',
  },
};
