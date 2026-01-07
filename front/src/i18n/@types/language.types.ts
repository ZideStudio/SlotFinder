/**
 * Language types for internationalization
 */

/**
 * Supported languages in the application
 */
export type Language = 'en' | 'fr';

/**
 * Language configuration object
 */
export type LanguageConfig = {
  code: Language;
  name: string;
  nativeName: string;
};

/**
 * Available language configurations
 */
export const SUPPORTED_LANGUAGES: Record<Language, LanguageConfig> = {
  en: {
    code: 'en',
    name: 'English',
    nativeName: 'English',
  },
  fr: {
    code: 'fr',
    name: 'French',
    nativeName: 'Français',
  },
} as const;

/**
 * Default language for the application
 */
export const DEFAULT_LANGUAGE: Language = 'en';

/**
 * Type guard to check if a string is a valid Language
 */
export const isValidLanguage = (lang: string): lang is Language => {
  return (Object.keys(SUPPORTED_LANGUAGES) as Language[]).includes(lang as Language);
};
