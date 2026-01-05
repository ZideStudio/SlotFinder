import ErrorIcon from "@material-symbols/svg-400/outlined/error.svg?react";
import './InputErrorMessage.scss';
import React from "react";
import { getClassName } from "@Front/utils/getClassName";

type InputErrorMessageProps = {
  children: React.ReactNode;
  className?: string;
  id?: string;
};

export const InputErrorMessage = ({ children, className, id }: InputErrorMessageProps) => {
  const isEmpty = !children || (typeof children === 'string' && children.trim() === '');

  if (isEmpty) {
    return null;
  }

  const parentClassName = getClassName({
    defaultClassName: 'ds-input-error-message',
    className,
  });

  return (
    <div className={parentClassName} role="alert">
      <ErrorIcon className="input-error-icon" aria-hidden="true" />
      <span id={id} className="input-error__message">
        {children}
      </span>
    </div>
  );
};