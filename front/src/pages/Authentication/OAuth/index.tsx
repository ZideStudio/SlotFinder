import { useTranslation } from 'react-i18next';
import { useOAuth } from './useOAuth';

import './index.css';

export const OAuth = () => {
  const { t } = useTranslation('authentication');
  const { oAuthProviders, errorCode } = useOAuth();

  return (
    <nav className="oauth-nav subgrid" aria-labelledby="oauth-provider-heading">
      <h2 id="oauth-provider-heading" style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '0.75rem' }}>
        {t('signInWithProvider')}
      </h2>
      <ul>
        {oAuthProviders.map(provider => (
          <li key={provider.label}>
            <a href={provider.href} aria-label={t(provider.ariaLabel)} rel="noopener noreferrer">
              {provider.icon}
              <span>{provider.label}</span>
            </a>
          </li>
        ))}
        {errorCode && (
          <span role="alert" style={{ color: 'red' }}>
            {t(`error.${errorCode}`)}
          </span>
        )}
      </ul>
    </nav>
  );
};
