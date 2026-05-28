import { Field } from "@Front/components/fields/Field/Field";
import { NumberInput } from "@Front/ui/molecules/Inputs/NumberInput/NumberInput";
import { type ComponentProps } from "react";

type NumberInputProps = Omit<ComponentProps<typeof NumberInput>, "error">;

export const NumberField = (props: NumberInputProps) => (
  <Field input={NumberInput} {...props} />
);
