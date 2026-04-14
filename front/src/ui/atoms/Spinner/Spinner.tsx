import { getClassName } from '@Front/utils/getClassName';
import './Spinner.scss';

type SpinnerProps = {
  className?: string;
  'aria-label'?: string;
};

export const Spinner = ({ className, 'aria-label': ariaLabel = "Chargement en cours" }: SpinnerProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-spinner',
    className,
  });

  return (
    <div className={parentClassName} role="status">
      <span className='ds-spinner__label'>{ariaLabel}</span>
    </div>
  );
};
