import { type ComponentProps, useState } from 'react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { ColorInput } from './ColorInput';
import { fn } from 'storybook/test';

const ColorInputStory = ({ onChange, ...args }: ComponentProps<typeof ColorInput>) => {
  const [value, setValue] = useState('');
  return (
    <ColorInput
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
  title: 'Molecules/Inputs/ColorInput',
  component: ColorInput,
  args: {
    label: 'Label',
    name: 'color-input',
    description: 'Choisir une couleur',
    value: '',
    onChange: fn(),
    required: false,
    className: 'custom-class',
  },
  argTypes: {
    onChange: { action: true, table: { disable: true } },
  },
  render: args => <ColorInputStory {...args} />,
} satisfies Meta<typeof ColorInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
