import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { useToastService } from '@Front/hooks/useToastService';
import { ToastProvider } from '@Front/providers/ToastProvider/ToastProvider';

const meta = {
  title: 'Atoms/Toast',
} satisfies Meta;

export default meta;

const ToastStoryContent = () => {
  const toastService = useToastService();

  return <button onClick={() => toastService.addToast('Ceci est un toast')}>Afficher le toast</button>;
};

export const Default: StoryObj = {
  render: () => (
    <ToastProvider>
      <ToastStoryContent />
    </ToastProvider>
  ),
};
