import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { CheckboxInputAtom } from './CheckboxInputAtom';

const meta = {
  title: 'Atoms/CheckboxInputAtom',
  component: CheckboxInputAtom,
  args: { className: 'custom-class' },
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
} satisfies Meta<typeof CheckboxInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};
