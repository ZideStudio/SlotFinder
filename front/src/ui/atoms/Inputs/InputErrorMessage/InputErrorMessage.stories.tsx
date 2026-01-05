import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { InputErrorMessage } from './InputErrorMessage';

const meta = {
  title: 'Atoms/Inputs/InputErrorMessage',
  component: InputErrorMessage,
  args: {
    id: 'error-message-1',
    children: 'Label text',
    className: 'custom-class',
  },
} satisfies Meta<typeof InputErrorMessage>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: 'test',
  },
};
