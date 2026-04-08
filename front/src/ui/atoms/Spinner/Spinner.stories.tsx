import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { Spinner } from './Spinner';

const meta = {
  title: 'Atoms/Spinner',
  component: Spinner,
  args: {
    className: 'custom-class',
    size: '24px',
  },
  decorators: [
    Story => (
      <div style={{ width: '50px', height: '50px' }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Spinner>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    size: '240px',
  },
};
