import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { ColorInput } from './ColorInput';

const meta = {
  component: ColorInput,
} satisfies Meta<typeof ColorInput>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { name: 'colorInput' },
};
