import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { InputErrorMessage } from './InputErrorMessage';

const meta = {
  title: 'Atoms/Inputs/InputErrorMessage',
  component: InputErrorMessage,
  args: {
    message: 'This field is required',
    id: 'error-message-1',
  },
  argTypes: {
    message: { control: 'text' },
    id: { control: 'text' },
  },
} satisfies Meta<typeof InputErrorMessage>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
