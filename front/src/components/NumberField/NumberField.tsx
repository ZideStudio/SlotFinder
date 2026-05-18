import { NumberInput } from '@Front/ui/molecules/Inputs/NumberInput/NumberInput';
import { type ComponentProps } from 'react';
import { useFormContext } from 'react-hook-form';

type NumberInputProps = Omit<ComponentProps<typeof NumberInput>, 'error'>;

export const NumberField = (props : NumberInputProps) => {
  const { register, formState } = useFormContext();
  const { errors } = formState;

  return (
    <NumberInput
      {...props}
      {...(register(props.name))}
      error={errors[props.name]?.message as string}
    />
  );
};
