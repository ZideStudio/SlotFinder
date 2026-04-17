import { FileUploadInputAtom } from '@Front/ui/atoms/Inputs/FileUploadInputAtom/FileUploadInputAtom';
import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { Field } from '@Front/ui/utils/components/Field/Field';
import { getClassName } from '@Front/utils/getClassName';
import { type ChangeEvent, type ComponentProps, useCallback, useEffect, useState } from 'react';
import './PictureUploadInput.scss';

type PictureUploadInputProps = ComponentProps<typeof FileUploadInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const PictureUploadInput = ({ onChange, className, ...props }: PictureUploadInputProps) => {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];

      setPreviewUrl(prev => {
        if (prev) URL.revokeObjectURL(prev);
        return file?.type.startsWith('image/') ? URL.createObjectURL(file) : null;
      });

      onChange?.(event);
    },
    [onChange],
  );

  useEffect(() => {
    return () => {
      setPreviewUrl(prev => {
        if (prev) URL.revokeObjectURL(prev);
        return null;
      });
    };
  }, []);

  const parentClassName = getClassName({
    defaultClassName: 'ds-picture-upload-input',
    className,
  });

  return (
    <div className={parentClassName}>
      <Field input={FileUploadInputAtom} onChange={handleChange} {...props} />
      {previewUrl && <img className="ds-picture-upload-input__preview" src={previewUrl} alt="Preview" />}
    </div>
  );
};
