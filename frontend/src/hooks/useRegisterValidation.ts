import { useState, useEffect } from 'react';
import { messages } from '../constants/messages';

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

  useEffect(() => {
    const newErrors: RegisterErrors = { email: '', password: '', repeatPassword: '', general: '' };
    if (touched.email) {
      if (!form.email) {
        newErrors.email = messages.emailRequired;
      } else if (!/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(form.email)) {
        newErrors.email = messages.emailInvalid;
      }
    }
    if (touched.password) {
      if (!form.password) {
        newErrors.password = messages.passwordRequired;
      } else if (form.password.length < 8) {
        newErrors.password = messages.passwordMinLength;
      } else if (!/[A-Z]/.test(form.password)) {
        newErrors.password = messages.passwordUppercase;
      } else if (!/[0-9]/.test(form.password)) {
        newErrors.password = messages.passwordNumber;
      } else if (!/[!@#$%^&*(),.?":{}|<>[\]/'_;+=-]/.test(form.password)) {
        newErrors.password = messages.passwordSymbol;
      }
    }
    if (touched.repeatPassword) {
      if (!form.repeatPassword) {
        newErrors.repeatPassword = messages.repeatPasswordRequired;
      } else if (form.repeatPassword !== form.password) {
        newErrors.repeatPassword = messages.passwordsDoNotMatch;
      }
    }
    setErrors(newErrors);
  }, [form, touched]);

  return errors;
}
