import { type ComponentPropsWithRef, useId, useState } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import PaletteIcon from '@material-symbols/svg-400/outlined/palette.svg?react';

import './ColorInputAtom.scss';

type ColorInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name' | 'type'> & {
  name: string;
  description: string;
};

export const ColorInputAtom = ({ id, className, description, onChange, ...props }: ColorInputAtomProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const DEFAULT_COLOR = '#ffffff';
  const [value, setValue] = useState('');
  const contentColor = getContrastTextColor(value ? value : DEFAULT_COLOR);

  const parentClassName = getClassName({
    defaultClassName: 'ds-color-input',
    className,
  });

  return (
    <label
      htmlFor={inputId}
      className={parentClassName}
      style={{ backgroundColor: value ?? DEFAULT_COLOR, color: contentColor }}
    >
      {' '}
      <PaletteIcon className="ds-color-input__icon" style={{ fill: contentColor }} aria-hidden="true" />
      <span className="ds-color-input__value">{value === '' ? description : value}</span>
      <input
        id={inputId}
        type="color"
        value={value}
        className="ds-color-input__input"
        {...props}
        onChange={event => {
          setValue(event.target.value);
          onChange?.(event);
        }}
      />
    </label>
  );
};
