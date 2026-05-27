import { Field } from "@Front/components/Field/Field";
import { TextareaInput } from "@Front/ui/molecules/Inputs/TextareaInput/TextareaInput";
import { type ComponentProps } from "react";

type TextareaInputProps = Omit<ComponentProps<typeof TextareaInput>, "error">;

export const TextareaField = (props: TextareaInputProps) => (
  <Field input={TextareaInput} {...props} />
);
