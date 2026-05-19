import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { SelectInputAtom } from '@Front/ui/atoms/Inputs/SelectInputAtom/SelectInputAtom';
import { Input } from '@Front/ui/utils/components/Input/Input';
import { type ComponentProps } from 'react';

type SelectInputProps = ComponentProps<typeof SelectInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const SelectInput = (props: SelectInputProps) => <Input input={SelectInputAtom} {...props} />;
