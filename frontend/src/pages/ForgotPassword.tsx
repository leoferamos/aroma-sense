import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import ErrorState from '../components/ErrorState';
import LoadingSpinner from '../components/LoadingSpinner';
import WordGrid from '../components/WordGrid';
import { useForgotPassword } from '../hooks/useForgotPassword';
import { usePasswordValidation } from '../hooks/usePasswordValidation';

const ForgotPassword: React.FC = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<'email' | 'code'>('email');
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [touched, setTouched] = useState({ password: false, confirmPassword: false });
  const { requestReset, confirmReset, loading, error, success } = useForgotPassword();
  const { validatePassword, validatePasswordConfirmation } = usePasswordValidation();

  const validateEmail = (value: string) => value.trim().length >= 3 && value.includes('@');

  // Use password validation
  const passwordError = validatePassword(newPassword, touched.password);
  const confirmError = validatePasswordConfirmation(newPassword, confirmPassword, touched.confirmPassword);

  const handleEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateEmail(email)) return;

    const success = await requestReset(email);
    if (success) {
      setStep('code');
    }
  };

  const handleCodeSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTouched({ password: true, confirmPassword: true });

    if (code.length !== 6 || passwordError || confirmError) {
      return;
    }

    const success = await confirmReset(email, code, newPassword);
    if (success) {
      setTimeout(() => navigate('/login'), 2000);
    }
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* Left side - gray background */}
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

      {/* Right side - white box */}
      <div className="w-full md:w-1/2 bg-white flex items-center justify-center px-4 py-8 md:px-0 md:py-0 relative">
        <div className="w-full max-w-md px-4 md:px-8 py-8 md:py-12 rounded-lg shadow-md">
          <div className="flex flex-col items-center mb-8">
            <img src="/logo.png" alt="Logo" className="h-16 md:h-20 mb-4" />
            <h2 className="text-2xl md:text-3xl font-medium text-center" style={{ fontFamily: 'Poppins, sans-serif' }}>
              {step === 'email' ? 'Forgot Password' : 'Reset Password'}
            </h2>
            <p className="text-sm text-gray-600 mt-2 text-center">
              {step === 'email' 
                ? 'Enter your email and we will send a 6-digit code to reset your password.' 
                : 'Enter the 6-digit code from your email and choose a new password.'}
            </p>
          </div>

          {step === 'email' ? (
            <form className="flex flex-col gap-6" onSubmit={handleEmailSubmit} noValidate>
              <InputField
                label="Email"
                name="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                onBlur={() => {}}
                required
                autoComplete="email"
              />
              <FormError message={(!email || validateEmail(email)) ? '' : 'Invalid email.'} />

              <button
                type="submit"
                className={`w-full text-white text-lg font-medium py-3 rounded-full mt-2 transition-colors ${
                  loading || !email ? 'bg-gray-300 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
                }`}
                disabled={loading || !email}
              >
                {loading ? <LoadingSpinner message="Sending..." /> : 'Send code'}
              </button>

              {error && <ErrorState message={error} />}
              {success && <div className="text-sm text-green-600">{success}</div>}
            </form>
          ) : (
            <form className="flex flex-col gap-6" onSubmit={handleCodeSubmit} noValidate>
              {/* Code input with 6 boxes */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Verification Code
                </label>
                <input
                  type="text"
                  inputMode="numeric"
                  pattern="[0-9]*"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  placeholder="000000"
                  className="w-full text-center text-3xl font-bold tracking-[0.5em] px-4 py-3 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                  autoComplete="off"
                  autoFocus
                />
                <FormError message={code && code.length !== 6 ? 'Code must be 6 digits.' : ''} />
              </div>

              <InputField
                label="New Password"
                name="newPassword"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                onBlur={() => setTouched({ ...touched, password: true })}
                required
                autoComplete="new-password"
              />
              <FormError message={passwordError} />

              <InputField
                label="Confirm Password"
                name="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                onBlur={() => setTouched({ ...touched, confirmPassword: true })}
                required
                autoComplete="new-password"
              />
              <FormError message={confirmError} />

              <button
                type="submit"
                className={`w-full text-white text-lg font-medium py-3 rounded-full mt-2 transition-colors ${
                  loading || code.length !== 6 || !!passwordError || !!confirmError || !newPassword || !confirmPassword
                    ? 'bg-gray-300 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
                }`}
                disabled={loading || code.length !== 6 || !!passwordError || !!confirmError || !newPassword || !confirmPassword}
              >
                {loading ? <LoadingSpinner message="Resetting..." /> : 'Reset password'}
              </button>

              {error && <ErrorState message={error} />}
              {success && (
                <div className="text-sm text-green-600 text-center font-medium">
                  {success} Redirecting to login...
                </div>
              )}

              <button
                type="button"
                onClick={() => {
                  setStep('email');
                  setCode('');
                  setNewPassword('');
                  setConfirmPassword('');
                  setTouched({ password: false, confirmPassword: false });
                }}
                className="text-sm text-blue-600 underline"
              >
                ← Resend code
              </button>
            </form>
          )}

          <div className="mt-6 text-gray-700 text-base text-center">
            <Link to="/login" className="underline text-blue-600">Back to login</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;
