import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router';

export const ErrorPage = () => {
  const { t } = useTranslation('error');
  const location = useLocation();

  return (
    <main>
      <h1>{t('title')}</h1>
      <p role="alert">{location.state?.message || t('unexpected')}</p>
    </main>
  );
};
