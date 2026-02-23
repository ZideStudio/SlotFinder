import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { ColorInput } from './ColorInput';

const meta = {
  title: 'Molecules/Inputs/ColorInput',
  component: ColorInput,
  args: { label: 'Label', name: 'color-input', description: 'Choisir une couleur', required: false, className: 'custom-class' },
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
} satisfies Meta<typeof ColorInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: 'An error occurred',
  },
};
