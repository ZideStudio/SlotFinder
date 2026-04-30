import { Heading } from '@Front/ui/atoms/Heading/Heading';
import { Button } from '@Front/ui/molecules/Button/Button';
import { ClickIcon } from '@Front/ui/molecules/ClickIcon/ClickIcon';
import Close from '@material-symbols/svg-400/rounded/close.svg?react';
import { type ComponentProps, type ReactNode } from 'react';

import './OverlayContent.scss';

type OverlayContentProps = {
  title: string;
  titleId?: string;
  children: ReactNode;
  closeButtonProps?: Partial<ComponentProps<typeof ClickIcon>>;
  primaryButtonProps: ComponentProps<typeof Button>;
  secondaryButtonProps?: ComponentProps<typeof Button>;
};

export const OverlayContent = ({
  title,
  titleId,
  children,
  closeButtonProps,
  primaryButtonProps,
  secondaryButtonProps,
}: OverlayContentProps) => (
  <>
    <header className="ds-overlay-content__header">
      <Heading level={1} id={titleId}>
        {title}
      </Heading>
      <ClickIcon
        aria-label="Fermer la fenêtre"
        className="ds-overlay-content__button--close"
        icon={Close}
        type="button"
        {...closeButtonProps}
      />
    </header>

    <div className="ds-overlay-content__body">{children}</div>

    <footer className="ds-overlay-content__footer">
      {secondaryButtonProps && <Button {...secondaryButtonProps} />}
      <Button {...primaryButtonProps} />
    </footer>
  </>
);

OverlayContent.displayName = 'OverlayContent';
