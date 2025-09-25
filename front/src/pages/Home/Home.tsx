import { useTranslation } from 'react-i18next';

export const Home = () => {
  const { t } = useTranslation('home');

  return (
    <main>
      <h1>{t('welcome')}</h1>
    </main>
  );
};
