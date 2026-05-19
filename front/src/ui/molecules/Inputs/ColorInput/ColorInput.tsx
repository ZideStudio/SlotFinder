import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { ColorInputAtom } from '@Front/ui/atoms/Inputs/ColorInputAtom/ColorInputAtom';
import { Input } from '@Front/ui/utils/components/Input/Input';
import { type ComponentProps } from 'react';

type ColorInputProps = ComponentProps<typeof ColorInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const ColorInput = (props: ColorInputProps) => <Input input={ColorInputAtom} {...props} />;
