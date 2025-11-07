import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import ErrorState from '../components/ErrorState';
import LoadingSpinner from '../components/LoadingSpinner';
import WordGrid from '../components/WordGrid';
import { useForgotPassword } from '../hooks/useForgotPassword';

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState('');
  const { requestPasswordReset, loading, error, success } = useForgotPassword();

  const validate = (value: string) => value.trim().length >= 3 && value.includes('@');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate(email)) {
      return;
    }

    await requestPasswordReset(email);
    setEmail('');
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
              Forgot Password
            </h2>
            <p className="text-sm text-gray-600 mt-2 text-center">Enter your email and we will send instructions to reset your password.</p>
          </div>

          <form className="flex flex-col gap-6" onSubmit={handleSubmit} noValidate>
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
            <FormError message={(!email || validate(email)) ? '' : 'Invalid email.'} />

            <button
              type="submit"
              className={`w-full text-white text-lg font-medium py-3 rounded-full mt-2 transition-colors ${
                loading || !email ? 'bg-gray-300 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700 cursor-pointer'
              }`}
              disabled={loading || !email}
            >
              {loading ? <LoadingSpinner message="Sending..." /> : 'Send instructions'}
            </button>

            {error && <ErrorState message={error} />}
            {success && <div className="text-sm text-green-600">{success}</div>}
          </form>

          <div className="mt-6 text-gray-700 text-base text-center">
            <Link to="/login" className="underline text-blue-600">Back to login</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;
