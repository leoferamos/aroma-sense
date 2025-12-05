import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import enCommon from './locales/en/common.json';
import ptCommon from './locales/pt/common.json';
import enAdmin from './locales/en/admin.json';
import ptAdmin from './locales/pt/admin.json';
import enErrors from './locales/en/errors.json';
import ptErrors from './locales/pt/errors.json';

const resources = {
  en: {
    common: enCommon,
    admin: enAdmin,
    errors: enErrors,
  },
  pt: {
    common: ptCommon,
    admin: ptAdmin,
    errors: ptErrors,
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'pt',
    debug: import.meta.env.DEV,

    interpolation: {
      escapeValue: false, // React already escapes
    },

    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage'],
    },
  });

export default i18n;