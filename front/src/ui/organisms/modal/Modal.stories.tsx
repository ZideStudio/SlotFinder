import { useModal } from '@Front/ui/utils/hooks/useModal';
import type { ComponentProps } from 'react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { fn } from 'storybook/test';
import { Modal } from './Modal';

const meta = {
  title: 'Organisms/Modal',
  component: Modal,
  args: {
    title: 'Modal title',
    children: "Modal's content",
    primaryButtonProps: { children: 'Action', onClick: fn() },
    secondaryButtonProps: { children: 'Action 2', onClick: fn(), variant: 'secondary' },
  },
  argTypes: {
    title: { control: 'text' },
    children: { control: 'text' },
    primaryButtonProps: { table: { disable: true } },
    secondaryButtonProps: { table: { disable: true } },
    ref: { table: { disable: true } },
  },
  render: (args: ComponentProps<typeof Modal>) => {
    // oxlint-disable-next-line react-hooks/rules-of-hooks
    const { modalRef, openModal } = useModal();

    return (
      <>
        <button type="button" onClick={openModal}>
          Open Modal
        </button>
        <Modal ref={modalRef} {...args} />
      </>
    );
  },
} satisfies Meta<typeof Modal>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithOnlyOneButton: Story = {
  args: {
    secondaryButtonProps: undefined,
  },
};
