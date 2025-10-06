import { useTranslation } from 'react-i18next';

export const Dashboard = () => {
  const { t } = useTranslation('dashboard');

  return <h1>{t('title')}</h1>;
};
