import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { NumberInputAtom } from './NumberInputAtom';

const meta = {
  title: 'Atoms/Inputs/NumberInputAtom',
  component: NumberInputAtom,
  args: { name: 'number-input', placeholder: 'Enter a number', className: 'custom-class' },
  argTypes: {
    'aria-invalid': {
      control: { type: 'boolean' },
    },
    onChange: { action: true, table: { disable: true } },
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof NumberInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {
  args: {
    min: 0,
    max: 100,
    step: 5,
  },
};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    'aria-invalid': 'true',
  },
};
