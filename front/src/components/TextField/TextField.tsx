import { TextInput } from '@Front/ui/molecules/Inputs/TextInput/TextInput';
import { useController, type UseControllerProps, type FieldValues } from 'react-hook-form';

type TextInputProps<Type extends FieldValues> = UseControllerProps<Type> & 
  Omit<React.ComponentProps<typeof TextInput>, 'value' | 'onChange' | 'onBlur' | 'error'>;

export const TextField = <Type extends FieldValues>({
  name,
  control,
  rules,
  defaultValue,
  ...props
}: TextInputProps<Type>) => {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
    rules,
    defaultValue,
  });

  return (
    <TextInput
      {...props}
      {...field}
      error={error?.message}
      required={props.required}
    />
  );
};