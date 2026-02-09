import type { Meta, StoryObj } from 'storybook-react-rsbuild';
import React from 'react';

import { Toast } from './Toast';

const meta = {
  title: 'Atoms/Toast',
  component: Toast,
} satisfies Meta<typeof Toast>;

export default meta;

export const Default: StoryObj<typeof meta> = {
  args: {
    children: 'Ceci est un toast',
    visible: false,
    onClose: () => {},
  },
  render: args => {
    const [visible, setVisible] = React.useState(args.visible ?? false);

    React.useEffect(() => {
      if (!visible) return;

      const timer = setTimeout(() => {
        setVisible(false);
      }, 3000);

      return () => clearTimeout(timer);
    }, [visible]);

    return (
      <div>
        <button onClick={() => setVisible(true)}>Afficher le toast</button>

        <Toast {...args} visible={visible} onClose={() => setVisible(false)}>
          {args.children}
        </Toast>
      </div>
    );
  },
};
