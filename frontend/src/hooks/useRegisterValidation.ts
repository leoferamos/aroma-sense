import { useState, useEffect } from 'react';
import { messages } from '../constants/messages';
import { usePasswordValidation } from './usePasswordValidation';

export interface RegisterForm {
  email: string;
  password: string;
  repeatPassword: string;
}

export interface RegisterErrors {
  email: string;
  password: string;
  repeatPassword: string;
  general: string;
}

export interface RegisterTouched {
  email: boolean;
  password: boolean;
  repeatPassword: boolean;
}

export function useRegisterValidation(form: RegisterForm, touched: RegisterTouched) {
  const [errors, setErrors] = useState<RegisterErrors>({
    email: '',
    password: '',
    repeatPassword: '',
    general: '',
  });
  
  const { validatePassword, validatePasswordConfirmation } = usePasswordValidation();

  useEffect(() => {
    const newErrors: RegisterErrors = { email: '', password: '', repeatPassword: '', general: '' };
    
    // Email validation
    if (touched.email) {
      if (!form.email) {
        newErrors.email = messages.emailRequired;
      } else if (!/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(form.email)) {
        newErrors.email = messages.emailInvalid;
      }
    }
    
    // Password validation
    newErrors.password = validatePassword(form.password, touched.password);
    
    // Password confirmation validation
    newErrors.repeatPassword = validatePasswordConfirmation(
      form.password,
      form.repeatPassword,
      touched.repeatPassword
    );
    
    setErrors(newErrors);
  }, [form, touched, validatePassword, validatePasswordConfirmation]);

  return errors;
}
