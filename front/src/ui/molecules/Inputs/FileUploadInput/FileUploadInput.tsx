import { getClassName } from '@Front/utils/getClassName';
import { LabelInput } from '../../../atoms/Inputs/LabelInput/LabelInput';

import './FileUploadInput.scss';

type FileUploadInputAtomProps = {
  className?: string;
  label: string;
  description?: string;
  helper?: string;
};

export const FileUploadInputAtom = ({ className, label, description, helper }: FileUploadInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-file-upload-input-atom',
    className,
  });

  return (
    <div className={parentClassName}>
      <LabelInput inputId="file-upload-input">{label}</LabelInput>
      {description && <p className="ds-file-upload-input-atom__description">{description}</p>}
      <input type="file" id="file-upload-input" className="ds-file-upload-input-atom__input" />
      {helper && <p className="ds-file-upload-input-atom__helper">{helper}</p>}
    </div>
  );
};
