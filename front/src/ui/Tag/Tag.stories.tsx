import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { Tag } from './Tag';

const meta = {
  title: 'Atoms/Tag',
  component: Tag,
  args: {
    children: 'Tag text',
    className: 'custom-class',
  },
  argTypes: {
    children: { control: 'text' },
  },
} satisfies Meta<typeof Tag>;

export default meta;

export const Default: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag>Default filled</Tag>
      <Tag appearance="outlined">Default outlined</Tag>
    </div>
  ),
};

export const Success: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="success">Success filled</Tag>
      <Tag variant="success" appearance="outlined">
        Success outlined
      </Tag>
    </div>
  ),
};

export const Warning: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="warning">Warning filled</Tag>
      <Tag variant="warning" appearance="outlined">
        Warning outlined
      </Tag>
    </div>
  ),
};

export const Error: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="error">Error filled</Tag>
      <Tag variant="error" appearance="outlined">
        Error outlined
      </Tag>
    </div>
  ),
};

export const Ellipsis: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag>Very very long filled text </Tag>
      <Tag appearance="outlined">Very long outlined text</Tag>
    </div>
  ),
};
