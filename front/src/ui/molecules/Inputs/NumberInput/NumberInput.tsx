import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { NumberInputAtom } from '@Front/ui/atoms/Inputs/NumberInputAtom/NumberInputAtom';
import { Input } from '@Front/ui/utils/components/Input/Input';
import { type ComponentProps } from 'react';

type NumberInputProps = ComponentProps<typeof NumberInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const NumberInput = (props: NumberInputProps) => <Input input={NumberInputAtom} {...props} />;
