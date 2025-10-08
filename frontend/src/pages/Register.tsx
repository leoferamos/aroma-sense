/**
 * Register page component.
 * Renders the registration form with validation and error handling.
 */

import React, { useState, useEffect } from 'react';
import InputField from '../components/InputField';
import FormError from '../components/FormError';

const Register: React.FC = () => {
  const [form, setForm] = useState({
    email: '',
    password: '',
    repeatPassword: '',
  });
  const [errors, setErrors] = useState({
    email: '',
    password: '',
    repeatPassword: '',
    general: '',
  });
  const [touched, setTouched] = useState({
    email: false,
    password: false,
    repeatPassword: false,
  });
  const [showPassword, setShowPassword] = useState(false);
  const [showRepeatPassword, setShowRepeatPassword] = useState(false);

  useEffect(() => {
    const newErrors = { email: '', password: '', repeatPassword: '', general: '' };
    if (touched.email) {
      if (!form.email) {
        newErrors.email = 'Email is required.';
      } else if (!/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(form.email)) {
        newErrors.email = 'Enter a valid email address.';
      }
    }
    if (touched.password) {
      if (!form.password) {
        newErrors.password = 'Password is required.';
      } else if (form.password.length < 8) {
        newErrors.password = 'Password must be at least 8 characters.';
      } else if (!/[A-Z]/.test(form.password)) {
        newErrors.password = 'Password must contain at least one uppercase letter.';
      } else if (!/[0-9]/.test(form.password)) {
        newErrors.password = 'Password must contain at least one number.';
      } else if (!/[!@#$%^&*(),.?":{}|<>\[\]\/\\'_;+=-]/.test(form.password)) {
        newErrors.password = 'Password must contain at least one symbol.';
      }
    }
    if (touched.repeatPassword) {
      if (!form.repeatPassword) {
        newErrors.repeatPassword = 'Repeat your password.';
      } else if (form.repeatPassword !== form.password) {
        newErrors.repeatPassword = 'Passwords do not match.';
      }
    }
    setErrors(newErrors);
  }, [form, touched]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    setTouched({ ...touched, [e.target.name]: true });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setTouched({ email: true, password: true, repeatPassword: true });
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* Left side - gray background */}
      <div className="hidden md:flex md:w-1/2 bg-gray-200 items-center justify-center relative">
  {/* Fragrance image: small part overlapping the white square */}
        <img
          src="/fragance.png"
          alt="Fragrance"
          className="absolute top-1/2 right-[-60px] -translate-y-1/2 w-[38vw] max-w-[420px] min-w-[180px] h-auto object-contain z-10"
        />
      </div>
      {/* Right side - white box */}
      <div className="w-full md:w-1/2 bg-white flex items-center justify-center px-4 py-8 md:px-0 md:py-0 relative">
        <div className="w-full max-w-md px-4 md:px-8 py-8 md:py-12 rounded-lg shadow-md">
          <div className="flex flex-col items-center mb-8">
            <img src="/logo.png" alt="Logo" className="h-16 md:h-20 mb-4" />
            <h2 className="text-2xl md:text-3xl font-medium text-center" style={{ fontFamily: 'Poppins, sans-serif' }}>
              Create Account
            </h2>
          </div>
          <form className="flex flex-col gap-6" onSubmit={handleSubmit} noValidate>
            <InputField
              label="Email"
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              onBlur={handleBlur}
              autoComplete="email"
              required
            />
            <FormError message={errors.email} />
            <InputField
              label="Password"
              type={showPassword ? 'text' : 'password'}
              name="password"
              value={form.password}
              onChange={handleChange}
              onBlur={handleBlur}
              autoComplete="new-password"
              required
              rightIcon={
                !showPassword ? (
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 3l18 18M3 12s3.5-6 9-6c2.1 0 4.1.7 5.7 1.8M21 12s-3.5 6-9 6c-2.1 0-4.1-.7-5.7-1.8" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                ) : (
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 12s3.5-6 9-6 9 6 9 6-3.5 6-9 6-9-6-9-6z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                )
              }
              onRightIconMouseDown={() => setShowPassword(true)}
              onRightIconMouseUp={() => setShowPassword(false)}
              onRightIconMouseLeave={() => setShowPassword(false)}
            />
            <FormError message={errors.password} />
            <span className="text-xs text-gray-500 mt-1">Use at least 8 characters, numbers and symbols</span>
            <InputField
              label="Repeat Password"
              type={showRepeatPassword ? 'text' : 'password'}
              name="repeatPassword"
              value={form.repeatPassword}
              onChange={handleChange}
              onBlur={handleBlur}
              autoComplete="new-password"
              required
              rightIcon={
                !showRepeatPassword ? (
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 3l18 18M3 12s3.5-6 9-6c2.1 0 4.1.7 5.7 1.8M21 12s-3.5 6-9 6c-2.1 0-4.1-.7-5.7-1.8" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                ) : (
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 12s3.5-6 9-6 9 6 9 6-3.5 6-9 6-9-6-9-6z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                )
              }
              onRightIconMouseDown={() => setShowRepeatPassword(true)}
              onRightIconMouseUp={() => setShowRepeatPassword(false)}
              onRightIconMouseLeave={() => setShowRepeatPassword(false)}
            />
            <FormError message={errors.repeatPassword} />
            <button
              type="submit"
              className="w-full bg-gray-300 text-white text-lg font-medium py-3 rounded-full mt-2 cursor-pointer"
              disabled={!form.email || !form.password || !form.repeatPassword}
            >
              Create Account
            </button>
            <FormError message={errors.general} />
          </form>
          <div className="mt-6 text-gray-700 text-base text-center">
            Already have an account? <span className="underline cursor-pointer">Login</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register;

