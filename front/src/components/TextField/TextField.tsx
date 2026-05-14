import { TextInput } from '@Front/ui/molecules/Inputs/TextInput/TextInput';
import { RegisterOptions, useFormContext } from 'react-hook-form';

type TextInputProps = Omit<React.ComponentProps<typeof TextInput>, 'value' | 'onChange' | 'onBlur' | 'error'> & {
  name: string;
  rules?: RegisterOptions;
};

export const TextField = ({ name, rules, id, ...props }: TextInputProps) => {
  const { register, formState } = useFormContext();
  const { errors } = formState;

  return (
    <TextInput
      {...props}
      {...(register(name, { ...rules }) as Partial<TextInputProps>)}
      id={id ?? name}
      name={name}
      error={errors[name]?.message as string}
    />
  );
};
