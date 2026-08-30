import { DurationInput } from "@Front/components/DurationInput/DurationInput";
import { UNITS } from "@Front/utils/constants/units";
import { get, useFormContext, type FieldError } from "react-hook-form";

type DurationFieldProps = {
  name: string;
  legend: string;
  required?: boolean;
  className?: string;
};

export const DurationField = ({
  name,
  legend,
  required,
  className,
}: DurationFieldProps) => {
  const {
    register,
    formState: { errors },
  } = useFormContext();

  const fields = UNITS.map((unit) => {
    const fieldName = `${name}.${unit}`;
    const fieldError: FieldError = get(errors, fieldName);

    return {
      unit,
      ...register(fieldName),
      error: fieldError?.message,
    };
  });

  return (
    <DurationInput
      legend={legend}
      required={required}
      className={className}
      fields={fields}
    />
  );
};
