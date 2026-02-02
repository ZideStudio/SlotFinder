import { InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { TextareaInputAtom } from '@Front/ui/atoms/Inputs/TextareaInputAtom/TextareaInputAtom';
import { Field } from '@Front/ui/utils/components/Field';
import { type ComponentProps } from 'react';

type TextareaInputProps = ComponentProps<typeof TextareaInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const TextareaInput = (props: TextareaInputProps) => {
  return <Field input={TextareaInputAtom} {...props} />;
};
