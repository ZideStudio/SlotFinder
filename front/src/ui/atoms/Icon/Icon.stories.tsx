import GoogleIcon from "@Front/assets/svg/google_icon.svg";
import type { Meta, StoryObj } from "storybook-react-rsbuild";
import { Icon } from "./Icon";

const meta = {
  title: "Atoms/Icon",
  component: Icon,
  args: {
    icon: GoogleIcon,
    className: "custom-class",
  },
  decorators: [
    (Story) => (
      <div style={{ width: "50px", height: "50px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Icon>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    icon: GoogleIcon,
  },
};
