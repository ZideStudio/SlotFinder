import { Field } from "@Front/components/Field/Field";
import { TextInput } from "@Front/ui/molecules/Inputs/TextInput/TextInput";
import { type ComponentProps } from "react";

type TextInputProps = Omit<ComponentProps<typeof TextInput>, "error">;

export const TextField = (props: TextInputProps) => (
  <Field input={TextInput} {...props} />
);
