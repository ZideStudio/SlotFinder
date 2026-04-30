import { getClassName } from '@Front/utils/getClassName';
import type { ComponentProps, ReactNode } from 'react';
import { OverlayContent } from '../OverlayContent/OverlayContent';

import './Popover.scss';
import { Card } from '@Front/ui/atoms/Card/Card';

type PopoverProps = {
  id: string;
  title: string;
  children: ReactNode;
  primaryButtonProps: ComponentProps<typeof OverlayContent>['primaryButtonProps'];
  secondaryButtonProps?: ComponentProps<typeof OverlayContent>['secondaryButtonProps'];
} & Omit<ComponentProps<'section'>, 'id' | 'children'>;

export const Popover = ({
  id,
  className,
  title,
  primaryButtonProps,
  secondaryButtonProps,
  children,
  ...props
}: PopoverProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-popover',
    className,
  });

  return (
    <Card as="section" popover="auto" id={id} className={parentClassName} {...props}>
      <OverlayContent
        title={title}
        primaryButtonProps={primaryButtonProps}
        secondaryButtonProps={secondaryButtonProps}
        closeButtonProps={{ popoverTarget: id, popoverTargetAction: 'hide' }}
      >
        {children}
      </OverlayContent>
    </Card>
  );
};
