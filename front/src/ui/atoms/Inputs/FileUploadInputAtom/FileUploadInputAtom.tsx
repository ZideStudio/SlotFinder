import { getClassName } from '@Front/utils/getClassName';
import { useId, type ComponentPropsWithRef } from 'react';
import UploadIcon from '@material-symbols/svg-400/outlined/upload_file.svg?react';

import './FileUploadInputAtom.scss';

type FileUploadInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name' | 'type'> & {
  name: string;
};

export const FileUploadInputAtom = ({ id, className, ...props }: FileUploadInputAtomProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;

  const parentClassName = getClassName({
    defaultClassName: 'ds-file-upload-input-atom',
    className,
  });

  return (
    <label htmlFor={inputId} className={parentClassName}>
      <UploadIcon aria-hidden="true" className="ds-file-upload-input-atom__icon" />
      <span className="ds-file-upload-input-atom__description">Déposer</span>
      <input id={inputId} type="file" {...props} className="ds-file-upload-input-atom__input" />
    </label>
  );
};
