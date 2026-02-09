import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { ToastProvider } from '@Front/providers/ToastProvider';
import { useToastContext } from '@Front/hooks/useToastContext';

const meta = {
  title: 'Atoms/Toast',
} satisfies Meta;

export default meta;

const ToastStoryContent = () => {
  const { show } = useToastContext();

  return <button onClick={() => show('Ceci est un toast')}>Afficher le toast</button>;
};

export const Default: StoryObj = {
  render: () => (
    <ToastProvider>
      <ToastStoryContent />
    </ToastProvider>
  ),
};
