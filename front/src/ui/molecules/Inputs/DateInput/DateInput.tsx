import { type InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { type LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { DateInputAtom } from "@Front/ui/atoms/Inputs/DateInputAtom/DateInputAtom";
import { Input } from "@Front/ui/utils/components/Input/Input";
import { type ComponentProps } from "react";

type DateInputProps = ComponentProps<typeof DateInputAtom> & {
  label: ComponentProps<typeof LabelInput>["children"];
  error?: ComponentProps<typeof InputErrorMessage>["children"];
};

export const DateInput = (props: DateInputProps) => (
  <Input input={DateInputAtom} {...props} />
);
