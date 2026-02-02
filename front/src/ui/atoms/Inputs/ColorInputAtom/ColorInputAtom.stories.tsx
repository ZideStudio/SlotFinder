import type { Meta, StoryObj } from 'storybook-react-rsbuild';

import { ColorInputAtom } from './ColorInputAtom';

const meta = {
  component: ColorInputAtom,
} satisfies Meta<typeof ColorInputAtom>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { name: 'colorInput', description: 'Choisir une couleur' },
};

export const Error: Story = {
  args: {
    name: 'ErrorColorInput',
    'aria-invalid': 'true',
    description: 'Erreur de couleur',
  },
};
