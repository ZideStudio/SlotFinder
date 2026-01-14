import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { TextInput } from './TextInput';

const meta = {
  title: 'Molecules/Inputs/TextInput',
  component: TextInput,
  args: { label: 'Label', name: 'text-input', required: false, className: 'custom-class', placeholder: 'Enter text' },
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
} satisfies Meta<typeof TextInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
