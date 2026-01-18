import { type ComponentPropsWithRef, useEffect, useId, useState } from 'react';
import './ColorInput.scss';
import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import PaletteIcon from '@material-symbols/svg-400/outlined/palette.svg?react';

export type ColorInputProps = Omit<ComponentPropsWithRef<'input'>, 'name' | 'type'> & {
  name: string;
};

export const ColorInput = ({ name, className, ...props }: ColorInputProps) => {
  const id = useId();
  const [value, setValue] = useState('#FF0000');
  const contentColor = getContrastTextColor(value);

  const parentClassName = getClassName({
    defaultClassName: 'ds-color-input',
    className,
  });

  return (
    <label htmlFor={id} className={parentClassName} style={{ backgroundColor: value, color: contentColor }}>
      <PaletteIcon className="ds-color-input__icon" style={{fill: contentColor}} aria-hidden="true" />
      <span className="ds-color-input__value">{value}</span>

      <input
        id={id}
        type="color"
        name={name}
        value={value}
        className="ds-color-input__native"
        {...props}
        onChange={e => {
          setValue(e.target.value);
        }}
      />
    </label>
  );
};
