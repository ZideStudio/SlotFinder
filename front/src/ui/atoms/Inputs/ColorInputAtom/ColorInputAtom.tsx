import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import PaletteIcon from '@material-symbols/svg-400/outlined/palette.svg?react';
import { type ChangeEventHandler, type ComponentPropsWithRef, useId } from 'react';

import './ColorInputAtom.scss';

const DEFAULT_COLOR = '#ffffff';

type ColorInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name' | 'type' | 'value' | 'onChange'> & {
  name: string;
  description: string;
  value: string;
  onChange: ChangeEventHandler<HTMLInputElement>;
};

export const ColorInputAtom = ({ id, className, description, onChange, value, ...props }: ColorInputAtomProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const resolvedValue = value || DEFAULT_COLOR;
  const contentColor = getContrastTextColor(resolvedValue);

  const parentClassName = getClassName({
    defaultClassName: 'ds-color-input',
    className,
  });

  return (
    <label
      htmlFor={inputId}
      className={parentClassName}
      style={{ backgroundColor: resolvedValue, color: contentColor }}
    >
      <PaletteIcon className="ds-color-input__icon" style={{ fill: contentColor }} aria-hidden="true" />
      <span className="ds-color-input__value">{value === '' ? description : value}</span>
      <input
        id={inputId}
        type="color"
        className="ds-color-input__input"
        {...props}
        value={resolvedValue}
        onChange={onChange}
      />
    </label>
  );
};
