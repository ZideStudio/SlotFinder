import { appRoutes } from '@Front/routing/appRoutes';
import { useTranslation } from 'react-i18next';
import { NavLink } from 'react-router';

export const Welcome = () => {
  const { t } = useTranslation('welcome');

  return (
    <>
      <h1>{t('title')}</h1>
      <NavLink to={appRoutes.signUp()}>Sign Up</NavLink>
    </>
  );
};
