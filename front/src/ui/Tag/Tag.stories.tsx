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

export const InDecision: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="inDecision">In decision</Tag>
      <Tag variant="inDecision" appearance="outlined">
        In decision
      </Tag>
    </div>
  ),
};

export const ComingSoon: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="comingSoon">Coming soon</Tag>
      <Tag variant="comingSoon" appearance="outlined">
        Coming soon
      </Tag>
    </div>
  ),
};

export const Finished: StoryObj<typeof meta> = {
  render: () => (
    <div style={{ display: 'flex', gap: 12 }}>
      <Tag variant="finished">Finished</Tag>
      <Tag variant="finished" appearance="outlined">
        Finished
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
