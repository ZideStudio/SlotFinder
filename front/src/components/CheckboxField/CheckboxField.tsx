import { MonoCheckboxInput } from '@Front/ui/molecules/Inputs/MonoCheckboxInput/MonoCheckboxInput';
import { type ComponentProps } from 'react';
import { useFormContext } from 'react-hook-form';

type CheckboxInputProps = Omit<ComponentProps<typeof MonoCheckboxInput>, 'error'>;

export const CheckboxField = (props : CheckboxInputProps) => {
  const { register, formState } = useFormContext();
  const { errors } = formState;

  return (
    <MonoCheckboxInput
      {...props}
      {...(register(props.name))}
      error={errors[props.name]?.message as string}
    />
  );
};
