import { Button } from '@Front/ui/molecules/Button/Button';
import { useModal } from '@Front/ui/utils/hooks/useModal/useModal';
import { type ComponentProps, useEffect } from 'react';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { fn } from 'storybook/test';
import { Modal } from './Modal';

type ModalStoryProps = ComponentProps<typeof Modal> & { openOnMount?: boolean };

const ModalStory = ({ openOnMount = false, ...args }: ModalStoryProps) => {
  const { modalRef, openModal } = useModal();

  useEffect(() => {
    if (openOnMount) {
      openModal();
    }
    // oxlint-disable-next-line eslint-plugin-react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <Button onClick={openModal} style={{ width: '300px' }}>
        Open Modal
      </Button>
      <Modal ref={modalRef} {...args} />
    </>
  );
};

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
    ref: { table: { disable: true } },
    primaryButtonProps: { table: { disable: true } },
    secondaryButtonProps: { table: { disable: true } },
  },
  render: (args: ComponentProps<typeof Modal>) => <ModalStory {...args} />,
} satisfies Meta<typeof Modal>;

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

export const OpenByDefault: Story = {
  render: (args: ComponentProps<typeof Modal>) => <ModalStory {...args} openOnMount />,
};
