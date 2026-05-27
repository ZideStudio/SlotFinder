import { getClassName } from "@Front/utils/getClassName";
import "./Spinner.scss";

type SpinnerProps = {
  className?: string;
  label?: string;
};

export const Spinner = ({
  className,
  label = "Chargement en cours",
}: SpinnerProps) => {
  const parentClassName = getClassName({
    defaultClassName: "ds-spinner",
    className,
  });

  return (
    // oxlint-disable-next-line jsx-a11y/prefer-tag-over-role
    <div className={parentClassName} role="status">
      <span className="ds-spinner__label">{label}</span>
    </div>
  );
};
