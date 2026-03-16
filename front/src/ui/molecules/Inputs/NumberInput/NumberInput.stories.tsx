import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { NumberInput } from './NumberInput';

const meta = {
  title: 'Molecules/Inputs/NumberInput',
  component: NumberInput,
  args: {
    label: 'Label',
    name: 'number-input',
    required: false,
    className: 'custom-class',
    placeholder: 'Enter a number',
  },
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
} satisfies Meta<typeof NumberInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {
  args: {
    min: 0,
    max: 100,
    step: 1,
  },
};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
