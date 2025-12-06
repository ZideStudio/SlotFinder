import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { LabelInput } from './LabelInput';

const meta = {
  title: 'Atoms/LabelInput',
  component: LabelInput,
  args: {
    inputId: 'input-id',
    children: 'Label text',
    className: 'custom-class',
    required: false,
  },
  argTypes: {
    children: { control: 'text' },
  },
} satisfies Meta<typeof LabelInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};
