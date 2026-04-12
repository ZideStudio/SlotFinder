import { Button } from '@Front/ui/molecules/Button/Button';
import { useToastService } from '@Front/ui/utils/toast/hooks/useToastService';
import { ToastProvider } from '@Front/ui/utils/toast/toastProvider/ToastProvider';
import type { Meta, StoryObj } from 'storybook-react-rsbuild';

type ToastStoryArgs = {
  label: string;
  duration: number | null;
};

const meta = {
  title: 'Atoms/Toast',
  args: {
    label: 'This is a toast',
    duration: 3000,
  },
  argTypes: {
    label: {
      control: { type: 'text' },
      description: 'Toast message content',
    },
    duration: {
      control: { type: 'select' },
      // oxlint-disable-next-line no-magic-numbers
      options: [null, 1000, 2000, 3000, 4000, 5000],
      description: 'Duration in milliseconds (1000 to 5000) or null for persistent toast',
    },
  },
} satisfies Meta<ToastStoryArgs>;

export default meta;

const ToastStoryContent = ({ label, duration }: ToastStoryArgs) => {
  const toastService = useToastService();

  return (
    <Button onClick={() => toastService.addToast(label, duration)} style={{ width: '300px' }}>
      Ajouter un toast
    </Button>
  );
};

export const Default: StoryObj<ToastStoryArgs> = {
  render: ({ label, duration }) => (
    <ToastProvider>
      <ToastStoryContent label={label} duration={duration} />
    </ToastProvider>
  ),
};
