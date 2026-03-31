import { getClassName } from '@Front/utils/getClassName';
import './Spinner.scss';

type SpinnerProps = {
  className?: string;
};

export const Spinner = ({ className }: SpinnerProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-spinner',
    className,
  });

  return (
    <div aria-label="Chargement en cours..." aria-live="polite" aria-busy="true">
      <div className={parentClassName} role="presentation" aria-hidden="true"/>
    </div>
  );
};
