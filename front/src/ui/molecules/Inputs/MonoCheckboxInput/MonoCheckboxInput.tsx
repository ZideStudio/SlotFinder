import { CheckboxInputAtom } from '@Front/ui/atoms/Inputs/CheckboxInputAtom/CheckboxInputAtom';
import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { Input } from '@Front/ui/utils/components/Input/Input';
import { type ComponentProps } from 'react';
import './MonoCheckboxInput.scss';

type MonoCheckboxInputProps = ComponentProps<typeof CheckboxInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const MonoCheckboxInput = (props: MonoCheckboxInputProps) => <Input input={CheckboxInputAtom} defaultClassName='ds-mono-checkbox-input' {...props} />;
