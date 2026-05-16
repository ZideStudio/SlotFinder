import { TextInput } from '@Front/ui/molecules/Inputs/TextInput/TextInput';
import { type ComponentProps } from 'react';
import { useFormContext } from 'react-hook-form';

type TextInputProps = Omit<ComponentProps<typeof TextInput>, 'error'>;

export const TextField = (props : TextInputProps) => {
  const { register, formState } = useFormContext();
  const { errors } = formState;

  return (
    <TextInput
      {...props}
      {...(register(props.name))}
      error={errors[props.name]?.message as string}
    />
  );
};
