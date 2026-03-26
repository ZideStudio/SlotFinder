import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { DateInput } from './DateInput';

const meta = {
  title: 'Molecules/Inputs/DateInput',
  component: DateInput,
  args: { label: 'Label', name: 'date-input', required: false, className: 'custom-class' },
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
} satisfies Meta<typeof DateInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
