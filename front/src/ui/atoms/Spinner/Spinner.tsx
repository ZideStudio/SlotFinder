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

  return <div className={parentClassName} />;
};
