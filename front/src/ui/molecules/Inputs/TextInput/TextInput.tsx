import { InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { TextInputAtom } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';
import { Field } from '@Front/ui/utils/components/Field';
import { type ComponentProps } from 'react';

type TextInputProps = ComponentProps<typeof TextInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const TextInput = (props: TextInputProps) => {
  return <Field input={TextInputAtom} {...props} />;
};
