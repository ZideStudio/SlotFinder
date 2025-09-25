# i18n Directory

This directory contains all resources and configuration for internationalization (i18n) in the project.

## Purpose

- Centralize all translation files and i18n configuration in one place.
- Enable type-safe and autocompleted translation keys in TypeScript.
- Make it easy to add new namespaces (pages/features) or update translations.
- Keep translation logic separated from business and UI code.

## Structure

```
src/i18n/
  index.ts                # i18n initialization and configuration
  @types/
    i18next.d.ts          # i18next type augmentation
    resources.ts          # TypeScript types for translation keys (auto-generated)
  locales/
    en/
      home.json           # English translations for Home page
      ...                 # Add more namespaces as needed
      index.ts            # (optional) re-exports for namespaces
```

## How it works

- Translation files are written in JSON, organized by namespace (one file per page/feature).
- The `i18next-resources-for-ts` script scans all translation files and generates TypeScript types in `@types/resources.ts`.
- The i18n config (`index.ts`) loads all namespaces and provides them to i18next.
- You get autocompletion and type safety for translation keys in your components.

## Usage in the App

- Import `src/i18n` in your app entry point (e.g., `main.ts` or `App.tsx`).
- Use the `useTranslation` hook from `react-i18next` in your components.

## Example

```tsx
import { useTranslation } from 'react-i18next';

const { t } = useTranslation('home');

t('welcome');
```

## Adding or Updating Translations

1. Add or edit a JSON file in `locales/en/` (e.g., `profile.json`).
2. Run the script to update types:

   ```sh
   npm run i18next-resources-for-ts
   ```

3. Use the new keys in your code with autocompletion and type safety.

## Notes

- The types in `@types/resources.ts` is auto-generated. Do not edit it manually.
- You can add more languages by creating new folders in `locales/` and updating the i18n config.

---

For more details, see the [react-i18next documentation](https://react.i18next.com/) and [i18next-resources-for-ts](https://www.npmjs.com/package/i18next-resources-for-ts).
