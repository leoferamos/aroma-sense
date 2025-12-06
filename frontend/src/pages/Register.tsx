/**
 * Register page component.
 * Renders the registration form with validation and error handling.
 */

import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import ErrorState from '../components/ErrorState';
import WordGrid from '../components/WordGrid';
import { useRegisterValidation } from '../hooks/useRegisterValidation';
import { useRegister } from '../hooks/useRegister';
import { useTranslation } from 'react-i18next';

const STORAGE_KEY = 'register_form_draft_v1';
const INITIAL_FORM = {
  email: '',
  password: '',
  repeatPassword: '',
};

const getStoredDraft = () => {
  if (typeof window === 'undefined') {
    return { form: INITIAL_FORM, agreeTerms: false, agreePrivacy: false };
  }

  const stored = localStorage.getItem(STORAGE_KEY);
  if (!stored) {
    return { form: INITIAL_FORM, agreeTerms: false, agreePrivacy: false };
  }

  try {
    const parsed = JSON.parse(stored);
    return {
      form: { ...INITIAL_FORM, ...(parsed?.form || {}) },
      agreeTerms: Boolean(parsed?.agreeTerms),
      agreePrivacy: Boolean(parsed?.agreePrivacy),
    };
  } catch {
    localStorage.removeItem(STORAGE_KEY);
    return { form: INITIAL_FORM, agreeTerms: false, agreePrivacy: false };
  }
};

const Register: React.FC = () => {
  const navigate = useNavigate();
  const draft = getStoredDraft();
  const [form, setForm] = useState(draft.form);
  const [showSuccessOverlay, setShowSuccessOverlay] = useState(false);
  const [touched, setTouched] = useState({
    email: false,
    password: false,
    repeatPassword: false,
  });
  const [showPassword, setShowPassword] = useState(false);
  const [showRepeatPassword, setShowRepeatPassword] = useState(false);
  const [agreeTerms, setAgreeTerms] = useState(draft.agreeTerms);
  const [touchedAgree, setTouchedAgree] = useState(false);
  const [agreePrivacy, setAgreePrivacy] = useState(draft.agreePrivacy);
  const [touchedPrivacy, setTouchedPrivacy] = useState(false);

  const errors = useRegisterValidation(form, touched);
  const { register, loading, error, success } = useRegister();
  const { t } = useTranslation('common');

  useEffect(() => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({ form, agreeTerms, agreePrivacy })
    );
  }, [form, agreeTerms, agreePrivacy]);

  useEffect(() => {
    if (!success) return;

    setShowSuccessOverlay(true);
    setForm(INITIAL_FORM);
    setAgreeTerms(false);
    setAgreePrivacy(false);
    if (typeof window !== 'undefined') {
      localStorage.removeItem(STORAGE_KEY);
    }

    const timer = setTimeout(() => {
      navigate('/login');
    }, 2000);

    return () => clearTimeout(timer);
  }, [success, navigate]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    setTouched({ ...touched, [e.target.name]: true });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setTouched({ email: true, password: true, repeatPassword: true });
    setTouchedAgree(true);
    setTouchedPrivacy(true);
    if (!agreeTerms || !agreePrivacy) {
      return;
    }

    if (!errors.email && !errors.password && !errors.repeatPassword && form.email && form.password && form.repeatPassword && agreeTerms && agreePrivacy) {
      register({ email: form.email, password: form.password });
    }
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* Left side - gray background */}
      <div className="hidden md:flex md:w-1/2 items-center justify-center relative" style={{ background: '#EAECEF' }}>
        {/* Background words grid */}
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
      {/* Right side - white box */}
      <div className="w-full md:w-1/2 bg-white flex items-center justify-center px-4 py-8 md:px-0 md:py-0 relative">
        <div className="w-full max-w-md px-4 md:px-8 py-8 md:py-12 rounded-lg shadow-md relative overflow-hidden">
          {showSuccessOverlay && (
            <div className="absolute inset-0 z-20 bg-white/95 backdrop-blur-sm flex flex-col items-center justify-center text-center px-6">
              <h3 className="text-xl font-semibold text-gray-900 mb-2">{t('auth.registrationSuccessful')}</h3>
              <p className="text-sm text-gray-600">{t('auth.redirectingToLogin')}</p>
            </div>
          )}
          <div className="flex flex-col items-center mb-8">
            <img src="/logo.png" alt="Logo" className="h-16 md:h-20 mb-4" />
            <h2 className="text-2xl md:text-3xl font-medium text-center" style={{ fontFamily: 'Poppins, sans-serif' }}>
              {t('auth.createAccount')}
            </h2>
          </div>
          <form className="flex flex-col gap-6" onSubmit={handleSubmit} noValidate>
            <InputField
              label={t('auth.email')}
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
              label={t('auth.password')}
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
            <span className="text-xs text-gray-500 mt-1">{t('auth.passwordHelper')}</span>
            <InputField
              label={t('auth.repeatPassword')}
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

            <div className="flex items-start gap-3">
              <input
                id="agreeTerms"
                type="checkbox"
                checked={agreeTerms}
                onChange={(e) => setAgreeTerms(e.target.checked)}
                onBlur={() => setTouchedAgree(true)}
                className="mt-1 w-4 h-4"
              />
              <label htmlFor="agreeTerms" className="text-sm text-gray-700">
                {t('auth.agreeTerms')}{' '}
                <Link to="/terms" className="underline text-blue-600">{t('auth.termsOfService')}</Link>
              </label>
            </div>
            <FormError message={touchedAgree && !agreeTerms ? t('auth.agreeTermsError') : ''} />

            <div className="flex items-start gap-3">
              <input
                id="agreePrivacy"
                type="checkbox"
                checked={agreePrivacy}
                onChange={(e) => setAgreePrivacy(e.target.checked)}
                onBlur={() => setTouchedPrivacy(true)}
                className="mt-1 w-4 h-4"
              />
              <label htmlFor="agreePrivacy" className="text-sm text-gray-700">
                {t('auth.agreePrivacy')}{' '}
                <Link to="/privacy" className="underline text-blue-600">{t('auth.privacyPolicy')}</Link>
              </label>
            </div>
            <FormError message={touchedPrivacy && !agreePrivacy ? t('auth.agreePrivacyError') : ''} />
            {errors.general || error ? <ErrorState message={errors.general || error} /> : null}
            {success && <span className="text-green-600 text-sm mt-2">{success}</span>}
            <button
              type="submit"
              className={`w-full text-white text-lg font-medium py-3 rounded-full mt-2 transition-colors ${loading || !form.email || !form.password || !form.repeatPassword || !agreeTerms || !agreePrivacy
                ? 'bg-gray-300 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
                }`}
              disabled={loading || !form.email || !form.password || !form.repeatPassword || !agreeTerms || !agreePrivacy}
            >
              {loading ? t('auth.registering') : t('auth.createAccount')}
            </button>
          </form>
          <div className="mt-6 text-gray-700 text-base text-center">
            {t('auth.alreadyHaveAccount')} <Link to="/login" className="underline">{t('auth.login')}</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register;