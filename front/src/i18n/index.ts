import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enAuthentication from './locales/en/authentication.json';
import enDashboard from './locales/en/dashboard.json';
import enError from './locales/en/error.json';
import enSignUp from './locales/en/signUp.json';
import enWelcome from './locales/en/welcome.json';

// oxlint-disable-next-line no-named-as-default-member
i18n.use(initReactI18next).init({
  resources: {
    en: {
      authentication: enAuthentication,
      dashboard: enDashboard,
      error: enError,
      signUp: enSignUp,
      welcome: enWelcome,
    },
  },
  lng: 'en',
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
});
