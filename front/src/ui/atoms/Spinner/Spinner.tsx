import { getClassName } from '@Front/utils/getClassName';
import './Spinner.scss';

type SpinnerProps = {
  className?: string;
  size?: string;
  'aria-label'?: string;
};

export const Spinner = ({ className, size, 'aria-label': ariaLabel = "Chargement en cours" }: SpinnerProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-spinner',
    className,
  });

  return (
    <div className={parentClassName} aria-label={ariaLabel} role="status" style={size ? { '--ds-spinner-size': size } as React.CSSProperties : {}} />
  );
};
