import authentication from '../locales/en/authentication.json';
import dashboard from '../locales/en/dashboard.json';
import error from '../locales/en/error.json';
import signUp from '../locales/en/signUp.json';
import welcome from '../locales/en/welcome.json';

const resources = {
  authentication,
  dashboard,
  error,
  signUp,
  welcome
} as const;

export default resources;
