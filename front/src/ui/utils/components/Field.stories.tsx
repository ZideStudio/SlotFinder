import type { Meta, StoryObj } from '@storybook/react';

import { Field } from './Field';
import { TextInputAtom } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';

const meta = {
  title: 'Utils/Components/Field',
  component: Field,
  args: {
    label: "Nom d'utilisateur",
    input: TextInputAtom,
    className: 'custom-class',
  },
  argTypes: {
    error: { control: 'text' },
    onChange: { action: true, table: { disable: true } },
    input: { table: { disable: true } },
  },
  decorators: [
    Story => (
      <div style={{ width: '300px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Field>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithError: Story = {
  args: {
    label: "Email",
    error: "Format d'email invalide",
    required: true,
  },
};
