import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { CheckboxInputAtom } from './CheckboxInputAtom';

const meta = {
  title: 'Atoms/CheckboxInputAtom',
  component: CheckboxInputAtom,
  args: {
    className: 'custom-class',
    id: 'checkbox-1',
    name: 'checkbox-group',
    disabled: false,
    required: false,
  },
  argTypes: {
    id: { control: { type: 'text' } },
    name: { control: { type: 'text' } },
    disabled: { control: { type: 'boolean' } },
    required: { control: { type: 'boolean' } },
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
} satisfies Meta<typeof CheckboxInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};
