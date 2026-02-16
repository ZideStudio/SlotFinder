import { CheckboxInputAtom } from '@Front/ui/atoms/Inputs/CheckboxInputAtom/CheckboxInputAtom';
import { InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { Field } from '@Front/ui/utils/components/Field';
import { type ComponentProps } from 'react';

type MonoCheckboxInputProps = ComponentProps<typeof CheckboxInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
  isCheckbox?: boolean;
};

export const MonoCheckboxInput = (props: MonoCheckboxInputProps) => {
  return <Field input={CheckboxInputAtom} isCheckbox {...props} />;
};
