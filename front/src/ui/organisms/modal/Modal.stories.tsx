import { useModal } from '@Front/ui/utils/hooks/useModal';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { Modal } from './Modal';

const meta = {
  title: 'Organisms/Modal',
  component: Modal,
  args: {
    title: 'Modal title',
    children: "Modal's content",
  },
  argTypes: {
    title: { control: 'text' },
    children: { control: 'text' },
    className: { table: { disable: true } },
  },
} satisfies Meta<typeof Modal>;

export default meta;

const ModalStory = () => {
  const { modalRef, openModal } = useModal();

  return (
    <>
      <button type="button" onClick={openModal}>
        Open Modal
      </button>
      <Modal ref={modalRef} title="Modal title">
        Modal's content
      </Modal>
    </>
  );
};

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => <ModalStory />,
};
