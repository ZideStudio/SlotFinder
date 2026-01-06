import type { Meta, StoryObj } from '@storybook/react';

import { CheckboxInputAtom } from './CheckboxInputAtom';

const meta = {
  title: 'Atoms/Inputs/CheckboxInputAtom',
  component: CheckboxInputAtom,
  args: {
    className: 'custom-class',
    name: 'checkbox-group',
  },
  argTypes: {
    'aria-invalid': {
      control: { type: 'boolean' },
    },
    onChange: {
      action: 'changed',
      table: { disable: true },
    },
  },
} satisfies Meta<typeof CheckboxInputAtom>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    'aria-invalid': true,
  },
};
