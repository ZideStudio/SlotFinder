import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { TextInputAtom } from "@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom";
import { getClassName } from "@Front/utils/getClassName";
import { ComponentProps, useId } from "react";

interface FieldProps {
    input: React.ComponentType<any>;
    id?: string;
    label?: ComponentProps<typeof LabelInput>['children'];
    error?: ComponentProps<typeof InputErrorMessage>['children'];
    defaultClassName?: string;
    required?: boolean;
    className?: string;
    [key: string]: any;
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
        {...props}
      />
      <InputErrorMessage id={errorId}>{error}</InputErrorMessage>
    </div>
  );
};
