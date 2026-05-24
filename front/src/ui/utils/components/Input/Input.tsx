import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { getClassName } from "@Front/utils/getClassName";
import {
  useId,
  type ComponentProps,
  type ElementType,
  type ReactNode,
} from "react";
import "./Input.scss";

type InputProps<ComponentType extends ElementType> = {
  input: ComponentType;
  id?: string;
  label?: ReactNode;
  error?: ReactNode;
  defaultClassName?: string;
  required?: boolean;
  className?: string;
} & Omit<
  ComponentProps<ComponentType>,
  "id" | "aria-describedby" | "aria-invalid" | "required"
>;

export const Input = <ComponentType extends ElementType>({
  input: InputComponent,
  id,
  error,
  label,
  required,
  className,
  defaultClassName = "ds-input",
  ...props
}: InputProps<ComponentType>) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const errorId = `${inputId}-error`;

  const parentClassName = getClassName({
    defaultClassName,
    className,
  });

  const Component = InputComponent as React.ElementType;

  return (
    <div className={parentClassName}>
      <LabelInput inputId={inputId} required={required}>
        {label}
      </LabelInput>
      <Component
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
