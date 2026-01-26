import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { getClassName } from "@Front/utils/getClassName";
import { ComponentProps, useId, ReactNode } from "react";
import './Field.scss';

type FieldProps = Omit<ComponentProps<'input'>, 'id'> & {
    input: React.ComponentType<any>;
    id?: string;
    label?: ReactNode;
    error?: ReactNode;
    defaultClassName?: string;
    required?: boolean;
    className?: string;
}

export const Field = ({ input: Input, id, error, label, required, className, defaultClassName = "ds-field", ...props }: FieldProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const errorId = `${inputId}-error`;

  const parentClassName = getClassName({
    defaultClassName,
    className,
  });

  return (
    <div className={parentClassName}>
      <LabelInput inputId={inputId} required={required}>
        {label}
      </LabelInput>
      <Input
        id={inputId}
        aria-describedby={error ? errorId : undefined}
        aria-invalid={Boolean(error)}
        required={required}
        {...props}
      />
      <InputErrorMessage id={errorId}>{error}</InputErrorMessage>
    </div>
  );
};
