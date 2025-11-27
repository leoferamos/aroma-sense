import { useCallback } from 'react';
import { messages } from '../constants/messages';

/**
 * Hook to validate password according to security rules.
 * Can be used in registration, password reset, and change password flows.
 */
export function usePasswordValidation() {
  /**
   * Validates a single password against security requirements.
   * @param password - The password to validate
   * @param touched - Whether the field has been touched/blurred
   * @returns Error message or empty string if valid
   */
  const validatePassword = useCallback((password: string, touched: boolean): string => {
    if (!touched) return '';
    
    if (!password) {
      return messages.passwordRequired;
    }
    if (password.length < 8) {
      return messages.passwordMinLength;
    }
    if (!/[A-Z]/.test(password)) {
      return messages.passwordUppercase;
    }
    if (!/[0-9]/.test(password)) {
      return messages.passwordNumber;
    }
    if (!/[!@#$%^&*(),.?":{}|<>[\]/'_;+=-]/.test(password)) {
      return messages.passwordSymbol;
    }
    
    return '';
  }, []);

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
      return messages.repeatPasswordRequired;
    }
    if (confirmPassword !== password) {
      return messages.passwordsDoNotMatch;
    }
    
    return '';
  }, []);

  return {
    validatePassword,
    validatePasswordConfirmation,
  };
}
