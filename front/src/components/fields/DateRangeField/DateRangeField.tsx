import { NumberField } from "@Front/components/fields/NumberField/NumberField";
import { getClassName } from "@Front/utils/getClassName";
import "./DateRangeField.scss";

type DateRangeFieldProps = {
  name: string;
  legend: string;
  required?: boolean;
  className?: string;
};

export const DateRangeField = ({
  legend,
  required,
  className,
}: DateRangeFieldProps) => {
  const parentClassName = getClassName({
    defaultClassName: "ds-date-range-field",
    className,
  });

  return (
    <fieldset className={parentClassName}>
      <legend className="ds-date-range-field__legend">
        {legend}
        {required && " *"}
      </legend>

      <div className="ds-date-range-field__inputs">
        <NumberField name="days" label="Jour(s)" min={0} />
        <NumberField name="hours" label="Heure(s)" min={0} max={23} />
        <NumberField name="minutes" label="Minute(s)" min={0} max={59} />
      </div>
    </fieldset>
  );
};
