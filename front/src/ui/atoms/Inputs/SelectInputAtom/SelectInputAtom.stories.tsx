import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { SelectInputAtom } from './SelectInputAtom';

const options = [
  { label: 'Option 1', value: '1' },
  { label: 'Option 2', value: '2' },
  { label: 'Option 3', value: '3' },
];

const meta = {
  title: 'Atoms/Inputs/SelectInputAtom',
  component: SelectInputAtom,
  args: {
    id: 'select-input',
    name: 'select-input',
    placeholder: 'Select an option',
    options,
    className: 'custom-class',
  },
  argTypes: {
    'aria-invalid': {
      control: { type: 'boolean' },
    },
    onChange: { action: true, table: { disable: true }},
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof SelectInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    'aria-invalid': 'true',
  },
};
