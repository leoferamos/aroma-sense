import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import enCommon from './locales/en/common.json';
import ptCommon from './locales/pt/common.json';
import enAdmin from './locales/en/admin.json';
import ptAdmin from './locales/pt/admin.json';
import enErrors from './locales/en/errors.json';
import ptErrors from './locales/pt/errors.json';
import enLegal from './locales/en/legal.json';
import ptLegal from './locales/pt/legal.json';
import enPrivacy from './locales/en/privacy.json';
import ptPrivacy from './locales/pt/privacy.json';

const resources = {
  en: {
    common: enCommon,
    admin: enAdmin,
    errors: enErrors,
    legal: enLegal,
    privacy: enPrivacy,
  },
  pt: {
    common: ptCommon,
    admin: ptAdmin,
    errors: ptErrors,
    legal: ptLegal,
    privacy: ptPrivacy,
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