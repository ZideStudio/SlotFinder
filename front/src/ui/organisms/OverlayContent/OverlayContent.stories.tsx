import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { fn } from 'storybook/test';

import { OverlayContent } from './OverlayContent';

const meta = {
  title: 'Organisms/OverlayContent',
  component: OverlayContent,
  args: {
    title: 'OverlayContent title',
    children: "OverlayContent's body",
    primaryButtonProps: { children: 'Confirm', onClick: fn() },
    secondaryButtonProps: { children: 'Cancel', onClick: fn(), variant: 'secondary' },
    closeButtonProps: { onClick: fn() },
  },
  argTypes: {
    title: { control: 'text' },
    titleId: { control: 'text' },
    children: { control: 'text' },
    closeButtonProps: { table: { disable: true } },
  },
} satisfies Meta<typeof OverlayContent>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithOnlyOneButton: Story = {
  args: {
    secondaryButtonProps: undefined,
  },
  argTypes: {
    secondaryButtonProps: { table: { disable: true } },
  },
};
