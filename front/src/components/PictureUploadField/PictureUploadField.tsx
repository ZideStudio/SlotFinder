import { Field } from "@Front/components/Field/Field";
import { PictureUploadInput } from "@Front/ui/molecules/Inputs/PictureUploadInput/PictureUploadInput";
import { type ComponentProps } from "react";

type PictureUploadInputProps = Omit<
  ComponentProps<typeof PictureUploadInput>,
  "error"
>;

export const PictureUploadField = (props: PictureUploadInputProps) => (
  <Field input={PictureUploadInput} {...props} />
);
