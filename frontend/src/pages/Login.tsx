/**
 * Login page component.
 * Renders the login form with validation and error handling.
 */

import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import WordGrid from '../components/WordGrid';
import { useLoginValidation } from '../hooks/useLoginValidation';
import { messages } from '../constants/messages';
import { useLogin } from '../hooks/useLogin';

const Login: React.FC = () => {
  const [form, setForm] = useState({ email: '', password: '' });
  const [showPassword, setShowPassword] = useState(false);
  const { errors, validateForm } = useLoginValidation();
  const { login, loading, error, user } = useLogin();
  const [generalError, setGeneralError] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const valid = validateForm(form);
    if (valid) {
      await login(form);
      if (error) setGeneralError(error);
      else setGeneralError("");
    }
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      <div className="hidden md:flex md:w-1/2 items-center justify-center relative" style={{ background: '#EAECEF' }}>
        <div className="absolute inset-0 pl-4 pr-6 flex items-center overflow-hidden z-10">
          <WordGrid />
        </div>
        <img
          src="/fragance.png"
          alt="Fragrance"
          className="frag-mid frag-xl absolute top-1/2 right-[-120px] w-[42vw] max-w-[560px] min-w-[220px] lg:w-[48vw] xl:w-[52vw] h-auto object-contain z-30"
          style={{ transform: 'translateY(-50%) rotate(-20deg)' }}
        />
      </div>

      <div className="w-full md:w-1/2 bg-white flex items-center justify-center px-4 py-8 md:px-0 md:py-0 relative">
        <div className="w-full max-w-md px-4 md:px-8 py-8 md:py-12 rounded-lg shadow-md">
          <div className="flex flex-col items-center mb-8">
            <img src="/logo.png" alt="Logo" className="h-16 md:h-20 mb-4" />
            <h2 className="text-2xl md:text-3xl font-medium text-center" style={{ fontFamily: 'Poppins, sans-serif' }}>
              Login
            </h2>
          </div>

          <form className="flex flex-col gap-6" onSubmit={handleSubmit} noValidate>
            <InputField
              label="Email"
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
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
              autoComplete="current-password"
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

            <button
              type="submit"
              className="w-full bg-gray-300 text-white text-lg font-medium py-3 rounded-full mt-2"
              disabled={loading}
            >
              {loading ? "Logging in..." : messages.login}
            </button>
            <FormError message={generalError} />
          </form>

          <div className="mt-6 text-gray-700 text-base text-center">
            {messages.dontHaveAccount} <Link to="/register" className="underline">{messages.createOne}</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
