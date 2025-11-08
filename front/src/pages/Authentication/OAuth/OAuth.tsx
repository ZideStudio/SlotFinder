import { Grid } from '@Front/components/Grid/Grid';
import { useTranslation } from 'react-i18next';
import './OAuth.css';
import { useOAuth } from './useOAuth';

export const OAuth = () => {
  const { t } = useTranslation('authentication');
  const { oAuthProviders } = useOAuth();

  return (
    <Grid
      component="nav"
      container
      colSpan={{ 'desktop-small': 4, tablet: 4, mobile: 4 }}
      colStart={{ 'desktop-small': 5, tablet: 3, mobile: 1 }}
      aria-labelledby="oauth-provider-heading"
      className="oauth-nav"
    >
      <h2
        id="oauth-provider-heading"
        style={{
          fontSize: '1.1rem',
          fontWeight: 600,
          marginBottom: '0.75rem',
        }}
      >
        {t('signInWithProvider')}
      </h2>
      <ul style={{ display: 'flex', flexDirection: 'column', gap: 15 }}>
        {oAuthProviders.map(provider => (
          <li key={provider.label}>
            <a href={provider.href} aria-label={t(provider.ariaLabel)} rel="noopener noreferrer">
              {provider.icon}
              <span>{provider.label}</span>
            </a>
          </li>
        ))}
      </ul>
    </Grid>
  );
};
