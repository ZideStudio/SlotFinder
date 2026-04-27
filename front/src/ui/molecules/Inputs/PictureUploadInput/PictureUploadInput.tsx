import { FileUploadInputAtom } from '@Front/ui/atoms/Inputs/FileUploadInputAtom/FileUploadInputAtom';
import { type InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { type LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import { Field } from '@Front/ui/utils/components/Field/Field';
import { getClassName } from '@Front/utils/getClassName';
import { type ChangeEvent, type ComponentProps, useCallback, useRef } from 'react';
import './PictureUploadInput.scss';

type PictureUploadInputProps = ComponentProps<typeof FileUploadInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
  previewText?: string;
  defaultPreviewUrl?: string;
};

export const PictureUploadInput = ({
  onChange,
  className,
  previewText,
  defaultPreviewUrl,
  ...props
}: PictureUploadInputProps) => {
  const imgRef = useRef<HTMLImageElement | null>(null);
  const urlRef = useRef<string | null>(null);

  /**
   * React 19 ref callback attached to the <img> element.
   *
   * - On mount: stores the node and applies defaultPreviewUrl if provided.
   * - Returns a cleanup function called by React on unmount, which revokes
   *   any pending object URL to prevent memory leaks. React 19 never calls
   *   this callback with null when a cleanup function is returned, so no
   *   null guard is needed.
   */

  const imgCallbackRef = useCallback(
    (node: HTMLImageElement) => {
      imgRef.current = node;
      if (defaultPreviewUrl) {
        node.src = defaultPreviewUrl;
        node.hidden = false;
      }
      return () => {
        if (urlRef.current) {
          URL.revokeObjectURL(urlRef.current);
          urlRef.current = null;
        }
      };
    },
    [defaultPreviewUrl],
  );

  /**
   * Handles file selection changes.
   *
   * Revokes the previous object URL (if any) before creating a new one to
   * avoid memory leaks. Updates the <img> node directly — no state update,
   * no re-render. Non-image files hide the preview.
   */

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];

      if (urlRef.current) {
        URL.revokeObjectURL(urlRef.current);
      }

      const newUrl = file && file.type.startsWith('image/') ? URL.createObjectURL(file) : null;
      urlRef.current = newUrl;

      const img = imgRef.current;
      if (img) {
        img.src = newUrl ?? '';
        img.hidden = !newUrl;
      }
      onChange?.(event);
    },
    [onChange],
  );

  const parentClassName = getClassName({
    defaultClassName: 'ds-picture-upload-input',
    className,
  });

  return (
    <div className={parentClassName}>
      <Field
        input={FileUploadInputAtom}
        onChange={handleChange}
        accept="image/jpeg,image/png"
        multiple={false}
        {...props}
      />
      <img ref={imgCallbackRef} className="ds-picture-upload-input__preview" alt={previewText || 'Preview'} hidden />
    </div>
  );
};
