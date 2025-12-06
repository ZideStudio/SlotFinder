import { getClassName } from '@Front/utils/getClassName';
import type { ReactNode } from 'react';
import './LabelInput.scss';

type LabelInputProps = {
  inputId: string;
  children: ReactNode;
  className?: string;
  required?: boolean;
};

export const LabelInput = ({ inputId, children, className, required }: LabelInputProps) => {
  const parentClassName = getClassName({
      defaultClassName: 'ds-label-input',
      className,
    });
    
  return (
    <label htmlFor={inputId} className={parentClassName}>
      {children}
      {required ? <span aria-hidden>*</span> : null}
    </label>
  );
};
