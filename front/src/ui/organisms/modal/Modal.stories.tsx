import { useModal } from '@Front/ui/utils/hooks/useModal';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { Modal } from './Modal';

const meta = {
  title: 'Organisms/Modal',
  component: Modal,
  args: {
    title: 'Modal title',
    children: "Modal's content",
    primaryButtonProps: { children: 'Action' },
    secondaryButtonProps: { children: 'Close' },
  },
  argTypes: {
    title: { control: 'text' },
    children: { control: 'text' },
    className: { table: { disable: true } },
    primaryButtonProps: { table: { disable: true } },
    secondaryButtonProps: { table: { disable: true } },
  },
} satisfies Meta<typeof Modal>;

export default meta;

const DefaultModalStory = () => {
  const { modalRef, openModal } = useModal();

  return (
    <>
      <button type="button" onClick={openModal}>
        Open Modal
      </button>
      <Modal
        ref={modalRef}
        title="Modal title"
        primaryButtonProps={{ children: 'Action' }}
        secondaryButtonProps={{ children: 'Close' }}
      >
        Modal's content
      </Modal>
    </>
  );
};

const WithOnlyOneButtonModalStory = () => {
  const { modalRef, openModal } = useModal();

  return (
    <>
      <button type="button" onClick={openModal}>
        Open Modal
      </button>
      <Modal ref={modalRef} title="Modal title" primaryButtonProps={{ children: 'Action' }}>
        Modal's content
      </Modal>
    </>
  );
}; 

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => <DefaultModalStory />,
};

export const WithOnlyOneButton: Story = {
  render: () => <WithOnlyOneButtonModalStory />,
};
