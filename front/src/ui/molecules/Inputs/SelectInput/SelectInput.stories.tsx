import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { SelectInput } from './SelectInput';

const options = [
  { label: 'Option 1', value: '1' },
  { label: 'Option 2', value: '2' },
  { label: 'Option 3', value: '3' },
];

const meta = {
  title: 'Molecules/Inputs/SelectInput',
  component: SelectInput,
  args: {
    label: 'Label',
    name: 'select-input',
    required: false,
    className: 'custom-class',
    placeholder: 'Select an option',
    options,
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
} satisfies Meta<typeof SelectInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
