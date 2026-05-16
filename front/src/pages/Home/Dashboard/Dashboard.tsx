import { Header } from '@Front/components/Layout/Header/Header';
import { useTranslation } from 'react-i18next';

export const Dashboard = () => {
  const { t } = useTranslation('dashboard');

  return (
    <>
      <Header />
      <h1>{t('title')}</h1>
    </>
  );
};
