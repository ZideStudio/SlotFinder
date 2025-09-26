import { useSignUp } from '@Front/hooks/api/useSignUp';
import type { SignUpFormType } from '@Front/types/Authentication/signUp.types';
import { yupResolver } from '@hookform/resolvers/yup';
import { FormProvider, useForm } from 'react-hook-form';
import { object, string } from 'yup';
import { PASSWORD_MIN_LENGTH, USERNAME_MIN_LENGTH } from './constants';

const schema = object({
  username: string()
    .required('Username is required')
    .min(USERNAME_MIN_LENGTH, `Username must be at least ${USERNAME_MIN_LENGTH} characters`),
  email: string().email('Invalid email').required('Email is required'),
  password: string()
    .required('Password is required')
    .min(PASSWORD_MIN_LENGTH, `Password must be at least ${PASSWORD_MIN_LENGTH} characters`)
    .matches(/[A-Za-z]/, 'Password must contain letters')
    .matches(/\d/, 'Password must contain numbers')
    .matches(/[^A-Za-z0-9]/, 'Password must contain symbols'),
});

export const SignUp = () => {
  const { signUp, isLoading, error } = useSignUp();
  const methods = useForm<SignUpFormType>({
    resolver: yupResolver(schema),
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
            Sign Up
          </legend>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="username">Username</label>
            <input
              id="username"
              type="text"
              autoComplete="username"
              aria-describedby={methods.formState.errors.username ? 'username-error' : undefined}
              {...methods.register('username')}
            />
            {methods.formState.errors.username && (
              <span id="username-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.username.message}
              </span>
            )}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              aria-describedby={methods.formState.errors.email ? 'email-error' : undefined}
              {...methods.register('email')}
            />
            {methods.formState.errors.email && (
              <span id="email-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.email.message}
              </span>
            )}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              autoComplete="new-password"
              aria-describedby={methods.formState.errors.password ? 'password-error' : undefined}
              {...methods.register('password')}
            />
            {methods.formState.errors.password && (
              <span id="password-error" role="alert" style={{ color: 'red', marginTop: 2 }}>
                {methods.formState.errors.password.message}
              </span>
            )}
          </div>
        </fieldset>
        {error && (
          <span role="alert" style={{ color: 'red' }}>
            {error}
          </span>
        )}
        <button type="submit" style={{ marginTop: '1rem' }} disabled={isLoading}>
          Sign Up
        </button>
      </form>
    </FormProvider>
  );
};
