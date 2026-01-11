import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { TextareaInputAtom } from './TextareaInputAtom';

const meta = {
  title: 'Atoms/Inputs/TextareaInputAtom',
  component: TextareaInputAtom,
  args: { name: 'textarea-input', className: 'custom-class' },
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
} satisfies Meta<typeof TextareaInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    'aria-invalid': 'true',
  },
};
