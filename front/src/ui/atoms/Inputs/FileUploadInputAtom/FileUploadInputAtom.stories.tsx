import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { FileUploadInputAtom } from './FileUploadInputAtom';

const meta = {
  title: 'Atoms/Inputs/FileUploadInputAtom',
  component: FileUploadInputAtom,
} satisfies Meta<typeof FileUploadInputAtom>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    name: 'file-upload',
    className: 'custom-class',
  },
};

export const Error: Story = {
  args: {
    name: 'file-upload',
    'aria-invalid': 'true',
    className: 'custom-class',
  },
};
