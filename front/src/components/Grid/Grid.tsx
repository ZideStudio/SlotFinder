import { clsx } from '@Front/utils/getClassName';
import { type ComponentProps, type CSSProperties, type ElementType, type PropsWithChildren } from 'react';
import { getGridItemToken } from './getGridItemToken';
import type { Breakpoint, ColSpan, ColStart } from './types';

import './_grid.scss';

type ExtendableComponent<ComponentType extends ElementType> = {
  component?: ComponentType;
  className?: string;
  style?: CSSProperties;
} & Omit<ComponentProps<ComponentType>, 'className' | 'style'>;

export type GridOwnProps = {
  container?: boolean;
  colSpan?: ColSpan | Partial<Record<Breakpoint, ColSpan>>;
  colStart?: ColStart | Partial<Record<Breakpoint, ColStart>>;
};

export type GridProps<ComponentType extends ElementType> = GridOwnProps & ExtendableComponent<ComponentType>;

export const Grid = <ComponentType extends ElementType = 'div'>({
  component,
  container,
  className,
  children,
  colSpan,
  colStart,
  style,
  ...rest
}: PropsWithChildren<GridProps<ComponentType>>) => {
  const Component = component ?? 'div';

  return (
    <Component
      className={clsx(container && 'grid', (colSpan || colStart) && 'grid-item', className)}
      style={{
        ...getGridItemToken('col', colSpan ?? {}),
        ...getGridItemToken('start', colStart ?? {}),
        ...style,
      }}
      {...rest}
    >
      {children}
    </Component>
  );
};

Grid.displayName = 'Grid';
