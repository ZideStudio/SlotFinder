import { CheckboxInputAtom } from '@Front/ui/atoms/Inputs/CheckboxInputAtom/CheckboxInputAtom';
import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { Field } from '@Front/ui/utils/components/Field';
import { type ComponentProps } from 'react';
import './MonoCheckboxInput.scss';

type MonoCheckboxInputProps = ComponentProps<typeof CheckboxInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const MonoCheckboxInput = (props: MonoCheckboxInputProps) => <Field input={CheckboxInputAtom} defaultClassName='ds-mono-checkbox-input' {...props} />;
