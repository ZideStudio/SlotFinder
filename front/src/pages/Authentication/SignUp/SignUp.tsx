import { useSignUp } from '@Front/hooks/api/useSignUp';
import type { SignUpFormType } from '@Front/types/Authentication/signUp/signUp.types';
import { yupResolver } from '@hookform/resolvers/yup';
import { FormProvider, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { getSchema } from './validation';

export const SignUp = () => {
  const { signUp, isLoading, errorCode } = useSignUp();
  const { t } = useTranslation('signUp');
  const methods = useForm<SignUpFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <FormProvider {...methods}>
      <form
        onSubmit={methods.handleSubmit(signUp)}
        style={{ maxWidth: 400, margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '1.5rem' }}
        aria-labelledby="signup-legend"
      >
        <fieldset style={{ border: 'none', padding: 0, margin: 0 }}>
          <legend id="signup-legend" style={{ fontWeight: 'bold', marginBottom: '1rem' }}>
            {t('title')}
          </legend>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="username">{t('username')}</label>
            <input
              id="username"
              type="text"
              autoComplete="username"
              aria-describedby={methods.formState.errors.username ? 'username-error' : undefined}
              {...methods.register('username', { required: true })}
            />
            {methods.formState.errors.username && (
              <span id="username-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.username.message}
              </span>
            )}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="email">{t('email')}</label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              aria-describedby={methods.formState.errors.email ? 'email-error' : undefined}
              {...methods.register('email', { required: true })}
            />
            {methods.formState.errors.email && (
              <span id="email-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.email.message}
              </span>
            )}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="password">{t('password')}</label>
            <input
              id="password"
              type="password"
              autoComplete="new-password"
              aria-describedby={methods.formState.errors.password ? 'password-error' : undefined}
              {...methods.register('password', { required: true })}
            />
            {methods.formState.errors.password && (
              <span id="password-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.password.message}
              </span>
            )}
          </div>
        </fieldset>
        {errorCode && (
          <span role="alert" style={{ color: 'red' }}>
            {t(`error.${errorCode}`)}
          </span>
        )}
        <button type="submit" style={{ marginTop: '1rem' }} disabled={isLoading}>
          {t('submit')}
        </button>
      </form>
    </FormProvider>
  );
};
