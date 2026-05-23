import { Field } from '@Front/components/Field/Field';
import { ColorInput } from '@Front/ui/molecules/Inputs/ColorInput/ColorInput';
import { type ComponentProps } from 'react';

type ColorInputProps = Omit<ComponentProps<typeof ColorInput>, 'error'>;

export const ColorField = (props: ColorInputProps) => <Field input={ColorInput} {...props} />;
