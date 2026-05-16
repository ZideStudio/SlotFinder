import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { TextareaInputAtom } from '@Front/ui/atoms/Inputs/TextareaInputAtom/TextareaInputAtom';
import { Input } from '@Front/ui/utils/components/Input/Input';
import { type ComponentProps } from 'react';

type TextareaInputProps = ComponentProps<typeof TextareaInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const TextareaInput = (props: TextareaInputProps) => <Input input={TextareaInputAtom} {...props} />;
