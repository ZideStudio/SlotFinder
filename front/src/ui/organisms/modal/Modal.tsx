// oxlint-disable react/jsx-props-no-spreading
import { getClassName } from '@Front/utils/getClassName';
import type { HTMLAttributes, ReactNode } from 'react';

type ModalProps = {
  title: string;
  children: ReactNode;
  className?: string;
} & HTMLAttributes<HTMLDivElement>;

export const Modal = ({ title, className, children, ...props }: ModalProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-modal',
    className,
  });

  return (
    <div className={parentClassName} {...props}>
      {title}
      {children}
    </div>
  );
};
