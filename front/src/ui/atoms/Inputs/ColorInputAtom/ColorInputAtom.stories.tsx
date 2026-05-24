import { type ComponentProps, useState } from 'react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { ColorInputAtom } from './ColorInputAtom';
import { fn } from 'storybook/test';

const ColorInputAtomStory = ({ onChange, ...args }: ComponentProps<typeof ColorInputAtom>) => {
  const [value, setValue] = useState('');
  return (
    <ColorInputAtom
      {...args}
      value={value}
      onChange={e => {
        setValue(e.target.value);
        onChange?.(e);
      }}
    />
  );
};

const meta = {
  title: 'Atoms/Inputs/ColorInputAtom',
  component: ColorInputAtom,
  args: {
    name: 'colorInput',
    value: '',
    onChange: fn(),
  },
  argTypes: {
    onChange: { action: true, table: { disable: true } },
  },
  render: args => <ColorInputAtomStory {...args} />,
} satisfies Meta<typeof ColorInputAtom>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { description: 'Choisir une couleur' },
};

export const Error: Story = {
  args: {
    'aria-invalid': 'true',
    description: 'Erreur de couleur',
  },
};
