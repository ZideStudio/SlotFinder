import { NumberField } from "@Front/components/fields/NumberField/NumberField";
import { getClassName } from "@Front/utils/getClassName";
import "./DateRangeField.scss";
import { useTranslation } from "react-i18next";

type DateRangeFieldProps = {
  name: string;
  legend: string;
  required?: boolean;
  className?: string;
};

export const DateRangeField = ({
  name,
  legend,
  required,
  className,
}: DateRangeFieldProps) => {
  const parentClassName = getClassName({
    defaultClassName: "ds-date-range-field",
    className,
  });

  const { t } = useTranslation("duration");

  return (
    <fieldset className={parentClassName}>
      <legend className="ds-date-range-field__legend">
        {legend}
        {Boolean(required) && <span aria-hidden>*</span>}
      </legend>

      <div className="ds-date-range-field__inputs">
        <NumberField
          name={`${name}.days`}
          label={t("days")}
          min={0}
          required={required}
        />
        <NumberField
          name={`${name}.hours`}
          label={t("hours")}
          min={0}
          max={23}
          required={required}
        />
        <NumberField
          name={`${name}.minutes`}
          label={t("minutes")}
          min={0}
          max={59}
          required={required}
        />
      </div>
    </fieldset>
  );
};
