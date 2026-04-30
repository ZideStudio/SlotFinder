import { Button } from '@Front/ui/molecules/Button/Button';
import { usePopover } from '@Front/ui/utils/hooks/usePopover/usePopover';
import { type ComponentProps, useEffect } from 'react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { fn } from 'storybook/test';
import { Popover } from './Popover';

type PopoverStoryProps = ComponentProps<typeof Popover> & { openOnMount?: boolean };

const PopoverStory = ({ openOnMount = false, ...args }: PopoverStoryProps) => {
  const { triggerProps, popoverProps } = usePopover();

  useEffect(() => {
    if (openOnMount) {
      document.getElementById(popoverProps.id)?.showPopover();
    }
    // oxlint-disable-next-line eslint-plugin-react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <Button {...triggerProps} style={{ ...triggerProps.style, width: '300px' }}>
        Open popover
      </Button>
      <Popover {...args} {...popoverProps} />
    </>
  );
};

const meta = {
  title: 'Organisms/Popover',
  component: Popover,
  args: {
    id: 'story-popover',
    title: 'Popover title',
    children: "I'm a popover",
    primaryButtonProps: { children: 'Confirm', onClick: fn() },
    secondaryButtonProps: { children: 'Cancel', onClick: fn(), variant: 'secondary' },
  },
  argTypes: {
    id: { table: { disable: true } },
    children: { control: 'text' },
    primaryButtonProps: { table: { disable: true } },
    secondaryButtonProps: { table: { disable: true } },
  },
  render: (args: ComponentProps<typeof Popover>) => <PopoverStory {...args} />,
} satisfies Meta<typeof Popover>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithOnlyPrimaryButton: Story = {
  args: {
    secondaryButtonProps: undefined,
  },
  argTypes: {
    secondaryButtonProps: { table: { disable: true } },
  },
};

export const OpenByDefault: Story = {
  render: (args: ComponentProps<typeof Popover>) => <PopoverStory {...args} openOnMount />,
};
