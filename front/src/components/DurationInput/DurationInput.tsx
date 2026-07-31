import { useId } from "react";
import { get, useFormContext, type FieldError } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { getClassName } from "@Front/utils/getClassName";
import "./DurationInput.scss";
import { LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { NumberInputAtom } from "@Front/ui/atoms/Inputs/NumberInputAtom/NumberInputAtom";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";

type DurationInputProps = {
  name: string;
  legend: string;
  required?: boolean;
  className?: string;
};

type DurationUnit = "days" | "hours" | "minutes";

const UNIT_LIMITS: Record<DurationUnit, { min: number; max: number }> = {
  days: { min: 0, max: 21 },
  hours: { min: 0, max: 23 },
  minutes: { min: 0, max: 59 },
};

const UNITS: DurationUnit[] = ["days", "hours", "minutes"];

export const DurationInput = ({
  name,
  legend,
  required,
  className,
}: DurationInputProps) => {
  const parentClassName = getClassName({
    defaultClassName: "duration-field",
    className,
  });

  const { t } = useTranslation("duration");
  const baseId = useId();

  const {
    register,
    formState: { errors },
  } = useFormContext();

  return (
    <fieldset className={parentClassName}>
      <legend className="duration-field__legend">
        {legend}
        {Boolean(required) && (
          <span className="duration-field__legend-required" aria-hidden>
            *
          </span>
        )}
      </legend>

      <div className="duration-field__inputs">
        {UNITS.map((unit) => {
          const fieldName = `${name}.${unit}`;
          const fieldError: FieldError | undefined = get(errors, fieldName);
          const inputId = `${baseId}-${unit}`;
          const errorId = `${inputId}-error`;

          return (
            <div className="duration-field__field" key={unit}>
              <LabelInput inputId={inputId}>{t(unit)}</LabelInput>

              <NumberInputAtom
                id={inputId}
                min={UNIT_LIMITS[unit].min}
                max={UNIT_LIMITS[unit].max}
                required={required}
                aria-invalid={Boolean(fieldError)}
                aria-describedby={fieldError ? errorId : undefined}
                {...register(fieldName)}
              />

              <InputErrorMessage id={errorId}>{fieldError?.message}</InputErrorMessage>
            </div>
          );
        })}
      </div>
    </fieldset>
  );
};
