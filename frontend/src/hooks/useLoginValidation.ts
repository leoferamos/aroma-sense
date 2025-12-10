import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export interface LoginFormData {
  email: string;
  password: string;
}

export interface LoginErrors {
  email: string;
  password: string;
  general: string;
}

export const useLoginValidation = () => {
  const [errors, setErrors] = useState<LoginErrors>({
    email: '',
    password: '',
    general: ''
  });
  const { t } = useTranslation('common');

  const validateEmail = (email: string): boolean => {
    if (!email.trim()) {
      setErrors(prev => ({ ...prev, email: t('auth.emailRequired') }));
      return false;
    }
    
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setErrors(prev => ({ ...prev, email: t('auth.emailInvalid') }));
      return false;
    }
    
    setErrors(prev => ({ ...prev, email: '' }));
    return true;
  };

  const validatePassword = (password: string): boolean => {
    if (!password) {
      setErrors(prev => ({ ...prev, password: t('auth.passwordRequired') }));
      return false;
    }
    
    setErrors(prev => ({ ...prev, password: '' }));
    return true;
  };

  const validateForm = (formData: LoginFormData): boolean => {
    const isEmailValid = validateEmail(formData.email);
    const isPasswordValid = validatePassword(formData.password);
    
    return isEmailValid && isPasswordValid;
  };

  const clearErrors = () => {
    setErrors({ email: '', password: '', general: '' });
  };

  return {
    errors,
    validateEmail,
    validatePassword,
    validateForm,
    clearErrors,
    setErrors
  };
};