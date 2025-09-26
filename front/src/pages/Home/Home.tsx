import { appRoutes } from '@Front/routing/appRoutes';
import { useTranslation } from 'react-i18next';
import { NavLink } from 'react-router';

export const Home = () => {
  const { t } = useTranslation('home');

  return (
    <main>
      <h1>{t('welcome')}</h1>
      <NavLink to={appRoutes.signUp}>Sign Up</NavLink>
    </main>
  );
};
