import { TextareaInputAtom } from '@Front/ui/atoms/Inputs/TextareaInputAtom/TextareaInputAtom';
import { getClassName } from '@Front/utils/getClassName';
import { useId, type ComponentProps } from 'react';

import { InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import './TextareaInput.scss';

type TextareaInputProps = ComponentProps<typeof TextareaInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const TextareaInput = ({ id, className, label, error, ...textareaInputAtomProps }: TextareaInputProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const errorId = `${inputId}-error`;

  const parentClassName = getClassName({
    defaultClassName: 'ds-textarea-input',
    className,
  });

  return (
    <div className={parentClassName}>
      <LabelInput inputId={inputId} required={textareaInputAtomProps.required}>
        {label}
      </LabelInput>
      <TextareaInputAtom
        id={inputId}
        aria-describedby={error ? errorId : undefined}
        aria-invalid={Boolean(error)}
        {...textareaInputAtomProps}
      />
      <InputErrorMessage id={errorId}>{error}</InputErrorMessage>
    </div>
  );
};
