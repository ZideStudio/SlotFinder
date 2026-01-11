import { TextInputAtom } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';
import { getClassName } from '@Front/utils/getClassName';
import { useId, type ComponentProps } from 'react';

import { InputErrorMessage } from '@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage';
import { LabelInput } from '@Front/ui/atoms/Inputs/LabelInput/LabelInput';
import './TextInput.scss';

type TextInputProps = ComponentProps<typeof TextInputAtom> & {
  label: ComponentProps<typeof LabelInput>['children'];
  error?: ComponentProps<typeof InputErrorMessage>['children'];
};

export const TextInput = ({ id, className, label, error, ...textInputAtomProps }: TextInputProps) => {
  const generatedId = useId();
  const inputId = id || generatedId;
  const errorId = `${inputId}-error`;

  const parentClassName = getClassName({
    defaultClassName: 'ds-text-input',
    className,
  });

  return (
    <div className={parentClassName}>
      <LabelInput inputId={inputId} required={textInputAtomProps.required}>
        {label}
      </LabelInput>
      <TextInputAtom
        id={inputId}
        aria-describedby={error ? errorId : undefined}
        aria-invalid={Boolean(error)}
        {...textInputAtomProps}
      />
      <InputErrorMessage id={errorId}>{error}</InputErrorMessage>
    </div>
  );
};
