import { type ComponentPropsWithRef, useId, useState } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import PaletteIcon from '@material-symbols/svg-400/outlined/palette.svg?react';

import './ColorInputAtom.scss';

type ColorInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name' | 'type'> & {
  name: string;
};

export const ColorInputAtom = ({ id, className, onChange, ...props }: ColorInputAtomProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const [value, setValue] = useState('#FF0000');
  const contentColor = getContrastTextColor(value);

  const parentClassName = getClassName({
    defaultClassName: 'ds-color-input',
    className,
  });

  return (
    <label htmlFor={inputId} className={parentClassName} style={{ backgroundColor: value, color: contentColor }}>
      <PaletteIcon className="ds-color-input__icon" style={{ fill: contentColor }} aria-hidden="true" />
      <span className="ds-color-input__value">{value}</span>

      <input
        id={inputId}
        type="color"
        value={value}
        className="ds-color-input__input"
        {...props}
        onChange={e => {
          setValue(e.target.value);
          onChange?.(e);
        }}
      />
    </label>
  );
};
