import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { useToast } from '@Front/hooks/useToast';
import { ToastProvider } from '@Front/providers/ToastProvider';

const meta = {
  title: 'Atoms/Toast',
} satisfies Meta;

export default meta;

const ToastStoryContent = () => {
  const { show } = useToast();

  return <button onClick={() => show('Ceci est un toast')}>Afficher le toast</button>;
};

export const Default: StoryObj = {
  render: () => (
    <ToastProvider>
      <ToastStoryContent />
    </ToastProvider>
  ),
};
