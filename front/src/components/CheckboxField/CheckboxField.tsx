import { Field } from "@Front/components/Field/Field";
import { MonoCheckboxInput } from "@Front/ui/molecules/Inputs/MonoCheckboxInput/MonoCheckboxInput";
import { type ComponentProps } from "react";

type CheckboxInputProps = Omit<
  ComponentProps<typeof MonoCheckboxInput>,
  "error"
>;

export const CheckboxField = (props: CheckboxInputProps) => (
  <Field input={MonoCheckboxInput} {...props} />
);
