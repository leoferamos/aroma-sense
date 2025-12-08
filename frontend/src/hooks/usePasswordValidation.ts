import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';

/**
 * Hook to validate password according to security rules.
 * Can be used in registration, password reset, and change password flows.
 */
export function usePasswordValidation() {
  const { t } = useTranslation('common');
  /**
   * Validates a single password against security requirements.
   * @param password - The password to validate
   * @param touched - Whether the field has been touched/blurred
   * @returns Error message or empty string if valid
   */
  const validatePassword = useCallback((password: string, touched: boolean): string => {
    if (!touched) return '';
    
    if (!password) {
      return t('auth.passwordRequired');
    }
    if (password.length < 8) {
      return t('auth.passwordMinLength');
    }
    if (!/[A-Z]/.test(password)) {
      return t('auth.passwordUppercase');
    }
    if (!/[0-9]/.test(password)) {
      return t('auth.passwordNumber');
    }
    if (!/[!@#$%^&*(),.?":{}|<>[\]/'_;+=-]/.test(password)) {
      return t('auth.passwordSymbol');
    }
    
    return '';
  }, [t]);

  /**
   * Validates password confirmation matches the original password.
   * @param password - The original password
   * @param confirmPassword - The confirmation password
   * @param touched - Whether the confirmation field has been touched
   * @returns Error message or empty string if valid
   */
  const validatePasswordConfirmation = useCallback((
    password: string,
    confirmPassword: string,
    touched: boolean
  ): string => {
    if (!touched) return '';
    
    if (!confirmPassword) {
      return t('auth.confirmPassword');
    }
    if (confirmPassword !== password) {
      return t('auth.passwordsDoNotMatch');
    }
    
    return '';
  }, [t]);

  return {
    validatePassword,
    validatePasswordConfirmation,
  };
}
