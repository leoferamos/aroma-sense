import React from 'react';
import { useTranslation } from 'react-i18next';

const LanguageSelector: React.FC = () => {
  const { i18n } = useTranslation('common');

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return (
    <select
      value={i18n.language}
      onChange={(e) => changeLanguage(e.target.value)}
      className="bg-white border border-gray-300 rounded px-2 py-1 text-sm"
    >
      <option value="pt">Português</option>
      <option value="en">English</option>
    </select>
  );
};

export default LanguageSelector;