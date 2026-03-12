import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { DateInputAtom } from './DateInputAtom';

const meta = {
  title: 'Atoms/Inputs/DateInputAtom',
  component: DateInputAtom,
  args: { name: 'date-input', className: 'custom-class', id: 'date-inputId' },
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
} satisfies Meta<typeof DateInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    'aria-invalid': 'true',
  },
};
