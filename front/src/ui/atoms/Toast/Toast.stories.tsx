import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import { ToastProvider } from '@Front/providers/ToastProvider';
import { useToast } from '../../../hooks/useToast';

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
