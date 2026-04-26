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
  primaryButtonProps: ComponentProps<typeof Button>;
  secondaryButtonProps?: ComponentProps<typeof Button>;
  closeOverlay: () => void;
};

export const OverlayContent = ({
  title,
  titleId,
  children,
  primaryButtonProps,
  secondaryButtonProps,
  closeOverlay,
}: OverlayContentProps) => (
  <>
    <header className="ds-overlay-content__header">
      <Heading level={1} id={titleId}>
        {title}
      </Heading>
      <ClickIcon
        aria-label="Fermer la fenêtre"
        onClick={closeOverlay}
        className="ds-overlay-content__button--close"
        icon={Close}
        type="button"
      />
    </header>

    <div className="ds-overlay-content__body">{children}</div>

    <footer className="ds-overlay-content__footer">
      {secondaryButtonProps && <Button type="button" {...secondaryButtonProps} />}
      <Button type="button" {...primaryButtonProps} />
    </footer>
  </>
);

OverlayContent.displayName = 'OverlayContent';
