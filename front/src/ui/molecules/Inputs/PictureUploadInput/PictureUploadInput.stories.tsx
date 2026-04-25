import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { PictureUploadInput } from './PictureUploadInput';

const meta = {
  title: 'Molecules/Inputs/PictureUploadInput',
  component: PictureUploadInput,
  args: { label: 'Label', name: 'picture-upload-input', required: false, className: 'custom-class' },
  argTypes: {
    onChange: { action: true, table: { disable: true } },
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof PictureUploadInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
