import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { Tag } from './Tag';

const meta = {
  title: 'Atoms/Tag',
  component: Tag,
  args: {
    children: 'Tag text',
    className: 'custom-class',
    title: 'Tag title',
  },
  argTypes: {
    children: { control: 'text' },
  },
} satisfies Meta<typeof Tag>;

export default meta;

export const Filled: StoryObj<typeof meta> = {
  args: {
    color: '#e3b0b0',
    appearance: 'filled',
    children: 'Filled Tag',
    title: 'Filled Tag',
  },
};

export const Outlined: StoryObj<typeof meta> = {
  args: {
    color: '#28a745',
    appearance: 'outlined',
    children: 'Outlined Tag',
    title: 'Outlined Tag',
  },
};

export const Ellipsis: StoryObj<typeof meta> = {
  args: {
    color: '#ff00ff',
    appearance: 'filled',
    children: 'Very very long filled text',
    title: 'Very very long filled text',
  },
};
