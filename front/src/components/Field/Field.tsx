import { createElement, type ComponentProps, type ElementType } from 'react';
import { get, useFormContext, type FieldError, type FieldValues, type Path } from 'react-hook-form';

type FieldProps<ComponentType extends ElementType, FormValues extends FieldValues = FieldValues> = {
  input: ComponentType;
  name: Path<FormValues>;
} & Omit<ComponentProps<ComponentType>, 'error' | 'name'>;

export const Field = <ComponentType extends ElementType, FormValues extends FieldValues = FieldValues>({
  input: Input,
  name,
  ...props
}: FieldProps<ComponentType, FormValues>) => {
  const {
    register,
    formState: { errors },
  } = useFormContext<FormValues>();
  const fieldError: FieldError | undefined = get(errors, name);

  return createElement(Input, { ...props, ...register(name), error: fieldError?.message });
};
