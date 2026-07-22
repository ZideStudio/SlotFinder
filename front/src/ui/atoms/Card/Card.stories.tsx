import type { Meta, StoryObj } from "storybook-react-rsbuild";

import { Card } from "./Card";

const meta = {
  title: "Atoms/Card",
  component: Card,
  args: {
    children: "Card",
    className: "custom-class",
    borderWeight: "default",
  },
  argTypes: {
    children: {
      control: { type: "text" },
    },
    borderWeight: {
      control: { type: "radio" },
      options: ["default", "bold"],
    },
  },
  decorators: [
    (Story) => (
      <div style={{ width: "300px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Card>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Bold: Story = {
  args: {
    borderWeight: "bold",
  },
};
