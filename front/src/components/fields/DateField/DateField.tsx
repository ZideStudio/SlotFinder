import { Field } from "@Front/components/fields/Field/Field";
import { DateInput } from "@Front/ui/molecules/Inputs/DateInput/DateInput";
import { type ComponentProps } from "react";

type DateInputProps = Omit<ComponentProps<typeof DateInput>, "error">;

export const DateField = (props: DateInputProps) => (
  <Field input={DateInput} {...props} />
);
