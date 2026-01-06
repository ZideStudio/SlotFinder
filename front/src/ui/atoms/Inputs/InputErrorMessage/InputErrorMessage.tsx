import ErrorIcon from '@material-symbols/svg-400/outlined/error.svg?react';
import './InputErrorMessage.scss';
import { ReactNode } from 'react';
import { getClassName } from '@Front/utils/getClassName';

type InputErrorMessageProps = {
  children: ReactNode;
  className?: string;
  id?: string;
};

export const InputErrorMessage = ({ children, className, id }: InputErrorMessageProps) => {
  if (!children) {
    return null;
  }

  const parentClassName = getClassName({
    defaultClassName: 'ds-input-error-message',
    className,
  });

  return (
    <div className={parentClassName} role="alert">
      <ErrorIcon className="ds-input-error-message__icon" aria-hidden="true" />
      <span id={id}>{children}</span>
    </div>
  );
};
