import { use } from "i18next";
import { initReactI18next } from "react-i18next";
import enAuthentication from "./locales/en/authentication.json";
import enDashboard from "./locales/en/dashboard.json";
import enDuration from "./locales/en/duration.json";
import enError from "./locales/en/error.json";
import enSignUp from "./locales/en/signUp.json";
import enWelcome from "./locales/en/welcome.json";
import enWhoAreYou from "./locales/en/whoAreYou.json";

// oxlint-disable-next-line vitest/require-hook, react-hooks/rules-of-hooks
use(initReactI18next).init({
  resources: {
    en: {
      authentication: enAuthentication,
      dashboard: enDashboard,
      error: enError,
      signUp: enSignUp,
      welcome: enWelcome,
      duration: enDuration,
      whoAreYou: enWhoAreYou,
    },
  },
  lng: "en",
  fallbackLng: "en",
  interpolation: { escapeValue: false },
});
