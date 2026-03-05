import { getClassName } from '@Front/utils/getClassName';
import './Card.scss';

type CardProps = {
  className?: string;
  children: React.ReactNode;
};

export const Card = ({ className, children }: CardProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-card',
    className,
  });

  return <div className={parentClassName}>{children}</div>;
};
