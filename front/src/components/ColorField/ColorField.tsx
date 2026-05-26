import { ColorInput } from '@Front/ui/molecules/Inputs/ColorInput/ColorInput';
import { type ComponentProps } from 'react';
import { Controller, useFormContext, type FieldValues, type Path } from 'react-hook-form';

type ColorFieldProps<FormValues extends FieldValues = FieldValues> = Omit<
  ComponentProps<typeof ColorInput>,
  'error' | 'value' | 'onChange' | 'onBlur' | 'name'
> & {
  name: Path<FormValues>;
};

export const ColorField = <FormValues extends FieldValues = FieldValues>({
  name,
  ...props
}: ColorFieldProps<FormValues>) => {
  const { control } = useFormContext<FormValues>();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState: { error } }) => (
        <ColorInput {...props} {...field} value={field.value ?? ''} error={error?.message} />
      )}
    />
  );
};
