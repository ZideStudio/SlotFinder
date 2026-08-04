import { type ComponentProps, useId } from "react";
import { useTranslation } from "react-i18next";
import { getClassName } from "@Front/utils/getClassName";
import "./DurationInput.scss";
import { LabelInput } from "@Front/ui/atoms/Inputs/LabelInput/LabelInput";
import { NumberInputAtom } from "@Front/ui/atoms/Inputs/NumberInputAtom/NumberInputAtom";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { type DurationUnit, UNIT_LIMITS } from "@Front/utils/units";

type UnitFieldProps = {
  unit: DurationUnit;
  error?: string;
} & ComponentProps<typeof NumberInputAtom>;

type DurationInputProps = {
  legend: string;
  required?: boolean;
  className?: string;
  fields: UnitFieldProps[];
};

export const DurationInput = ({
  legend,
  required,
  className,
  fields,
}: DurationInputProps) => {
  const parentClassName = getClassName({
    defaultClassName: "duration-field",
    className,
  });

  const { t } = useTranslation("duration");
  const baseId = useId();

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
        {fields.map(({ unit, error, ...inputProps }) => {
          const inputId = `${baseId}-${unit}`;
          const errorId = `${inputId}-error`;

          return (
            <div className="duration-field__field" key={unit}>
              <div className="duration-field__field-label">
                <NumberInputAtom
                  id={inputId}
                  min={UNIT_LIMITS[unit].min}
                  max={UNIT_LIMITS[unit].max}
                  required={required}
                  aria-invalid={Boolean(error)}
                  aria-describedby={error ? errorId : undefined}
                  {...inputProps}
                />
                <LabelInput inputId={inputId}>{t(unit)}</LabelInput>
              </div>

              <InputErrorMessage id={errorId}>{error}</InputErrorMessage>
            </div>
          );
        })}
      </div>
    </fieldset>
  );
};
