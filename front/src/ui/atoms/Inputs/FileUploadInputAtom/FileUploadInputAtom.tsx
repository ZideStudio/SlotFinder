import { getClassName } from '@Front/utils/getClassName';

import './FileUploadInputAtom.scss';
import { ComponentPropsWithRef } from 'react';

type FileUploadInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name'> & {
  name: string;
};

export const FileUploadInputAtom = ({ className }: FileUploadInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-file-upload-input-atom',
    className,
  });

  return <input type="file" className={parentClassName} />;
};
