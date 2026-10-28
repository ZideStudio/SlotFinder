import { useTranslation } from "react-i18next";
import { useLocation } from "react-router";
import { appRoutes } from "@Front/routing/appRoutes";
import { Link } from "@Front/ui/atoms/Link/Link";
import "./ExistingAccount.scss";

type ExistingAccountProps = Record<string, never>;

const LOGIN_PATH = "/login";

export const ExistingAccount = (_props: ExistingAccountProps = {}) => {
  const { pathname } = useLocation();
  // TODO: Select the "signUp" or "signIn" translation namespace based on the current page once sign-in is implemented.
  const { t } = useTranslation("signUp");
  const isSignUpPage = pathname === appRoutes.signUp();
  // TODO: Use appRoutes.login() when login route is added to appRoutes
  const isLoginPage = pathname === LOGIN_PATH;

  if (!isSignUpPage && !isLoginPage) {
    return null;
  }

  //TODO: Use noAccount and signUpLink translations on i18n login.json
  // const prompt = isSignUpPage ? "existingAccount" : "noAccount";
  // const linkLabel = isSignUpPage ? "login" : "signUpLink";
  const prompt = "existingAccount";
  const linkLabel = "login";
  const href = isSignUpPage ? LOGIN_PATH : appRoutes.signUp();

  return (
    <p className="authentication-account-prompt">
      {t(prompt)}{" "}
      <Link as="a" className="authentication-account-prompt__link" href={href}>
        {t(linkLabel)}
      </Link>
    </p>
  );
};
