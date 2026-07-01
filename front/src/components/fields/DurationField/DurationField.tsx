import { NumberField } from "@Front/components/fields/NumberField/NumberField";
import { getClassName } from "@Front/utils/getClassName";
import "./DurationField.scss";
import { useTranslation } from "react-i18next";

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
  const parentClassName = getClassName({
    defaultClassName: "duration-field",
    className,
  });

  const { t } = useTranslation("duration");

  return (
    <fieldset className={parentClassName}>
      <legend className={`${parentClassName}__legend`}>
        {legend}
        {Boolean(required) && (
          <span className={`${parentClassName}__legend-required`} aria-hidden>
            *
          </span>
        )}
      </legend>

      <div className={`${parentClassName}__inputs`}>
        <NumberField
          name={`${name}.days`}
          label={t("days")}
          min={0}
          max={21}
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
