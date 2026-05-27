import type { Meta, StoryObj } from "@storybook/react";

import { Input } from "./Input";
import { TextInputAtom } from "@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom";

const meta = {
  title: "Utils/Components/Input",
  component: Input,
  args: {
    label: "Nom d'utilisateur",
    input: TextInputAtom,
    className: "custom-class",
    required: false,
  },
  argTypes: {
    error: { control: "text" },
    onChange: { action: true, table: { disable: true } },
    input: { table: { disable: true } },
  },
  decorators: [
    (Story) => (
      <div style={{ width: "300px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithError: Story = {
  args: {
    label: "Email",
    error: "Format d'email invalide",
    required: true,
  },
};
