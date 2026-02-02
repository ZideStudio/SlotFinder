import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { TextareaInput } from './TextareaInput';

const meta = {
  title: 'Molecules/Inputs/TextareaInput',
  component: TextareaInput,
  args: { label: 'Label', name: 'textarea-input', required: false, className: 'custom-class', placeholder: 'Enter text' },
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
} satisfies Meta<typeof TextareaInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
