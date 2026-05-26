import { Field } from '@Front/components/Field/Field';
import { SelectInput } from '@Front/ui/molecules/Inputs/SelectInput/SelectInput';
import { type ComponentProps } from 'react';

type SelectInputProps = Omit<ComponentProps<typeof SelectInput>, 'error'>;

export const SelectField = (props: SelectInputProps) => <Field input={SelectInput} {...props} />;
