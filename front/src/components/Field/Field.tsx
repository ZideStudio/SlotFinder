import { type ComponentProps, type ElementType } from 'react';
import { get, type FieldValues, type Path, useFormContext } from 'react-hook-form';

type FieldProps<ComponentType extends ElementType, FormValues extends FieldValues = FieldValues> = {
  input: ComponentType;
  name: Path<FormValues>;
} & Omit<ComponentProps<ComponentType>, 'error' | 'name'>;

export const Field = <ComponentType extends ElementType, FormValues extends FieldValues = FieldValues>({
  input: Input,
  name,
  ...props
}: FieldProps<ComponentType, FormValues>) => {
  const { register, formState } = useFormContext<FormValues>();
  const { errors } = formState;
  const error = get(errors, name)?.message;
  const message = typeof error === 'string' ? error : undefined;

  const Component = Input as ElementType;

  return <Component {...props} {...register(name)} error={message} />;
};
