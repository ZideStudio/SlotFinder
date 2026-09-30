import type { Meta, StoryObj } from "storybook-react-rsbuild";

import { PictureUploadInput } from "./PictureUploadInput";

const meta = {
  title: "Molecules/Inputs/PictureUploadInput",
  component: PictureUploadInput,
  args: {
    label: "Label",
    name: "picture-upload-input",
    required: false,
    previewText: "Preview",
    previewVariant: "default",
    defaultPreviewUrl:
      "https://avatars.githubusercontent.com/u/152389914?s=60&v=4",
    className: "custom-class",
  },
  argTypes: {
    onChange: { action: true, table: { disable: true } },
  },
  decorators: [
    (Story) => (
      <div style={{ width: "300px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof PictureUploadInput>;

export default meta;

export const Default: StoryObj<typeof meta> = {};

export const DefaultWithoutDefaultPreview: StoryObj<typeof meta> = {
  args: {
    defaultPreviewUrl: undefined,
  },
};

export const WithRoundedPreview: StoryObj<typeof meta> = {
  args: {
    previewVariant: "rounded",
  },
};

export const Invalid: StoryObj<typeof meta> = {
  args: {
    error: "An error occurred",
  },
};
