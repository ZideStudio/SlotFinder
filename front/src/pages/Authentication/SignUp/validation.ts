import type { TFunction } from 'i18next';
import { object, ref, string } from 'yup';
import { EMAIL_REGEX, PASSWORD_MIN_LENGTH, PASSWORD_REGEX, USERNAME_MIN_LENGTH } from './constants';

export const getSchema = (translate: TFunction) =>
  object({
    username: string()
      .required(translate('requiredUsername'))
      .min(USERNAME_MIN_LENGTH, translate('minLengthUsername', { min: USERNAME_MIN_LENGTH })),
    email: string().required(translate('requiredEmail')).matches(EMAIL_REGEX, translate('invalidEmail')),
    password: string()
      .required(translate('requiredPassword'))
      .min(PASSWORD_MIN_LENGTH, translate('minLengthPassword', { min: PASSWORD_MIN_LENGTH }))
      .matches(PASSWORD_REGEX, translate('passwordComplexity')),
    confirmPassword: string()
      .required(translate('requiredConfirmPassword'))
      .oneOf([ref('password')], translate('passwordsDoNotMatch')),
  });
