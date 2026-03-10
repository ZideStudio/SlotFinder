import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { FileUploadInputAtom } from './FileUploadInputAtom';

const meta = {
  component: FileUploadInputAtom,
} satisfies Meta<typeof FileUploadInputAtom>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    name: 'file-upload',
  },
};
