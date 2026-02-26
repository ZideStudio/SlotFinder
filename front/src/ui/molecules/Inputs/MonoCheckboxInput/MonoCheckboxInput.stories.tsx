import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { MonoCheckboxInput } from './MonoCheckboxInput';

const meta = {
  title: 'Molecules/Inputs/MonoCheckboxInput',
  component: MonoCheckboxInput,
  args: { label: 'Label', name: 'checkbox-input', required: false, className: 'custom-class' },
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
} satisfies Meta<typeof MonoCheckboxInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
