import ErrorIcon from "@material-symbols/svg-400/outlined/error.svg?react";
import './InputErrorMessage.scss';

type InputErrorMessageProps = {
  message: string;
  id?: string;
};

export const InputErrorMessage = ({ message, id }: InputErrorMessageProps) => (
  <div className="input-error">
    <ErrorIcon className="input-error__icon" aria-hidden="true" />
    <span id={id} aria-live="assertive" className="input-error__message">
      {message}
    </span>
  </div>
);