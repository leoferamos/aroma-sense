import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import LoadingSpinner from '../components/LoadingSpinner';
import { changePassword } from '../services/profile';
import { isAxiosError } from 'axios';

const ChangePassword: React.FC = () => {
    const navigate = useNavigate();
    const [form, setForm] = useState({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
    });
    const [touched, setTouched] = useState({
        currentPassword: false,
        newPassword: false,
        confirmPassword: false,
    });
    const [showCurrentPassword, setShowCurrentPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
        setTouched({ ...touched, [e.target.name]: true });
    };

    const validateForm = () => {
        const errors: Record<string, string> = {};

        if (touched.currentPassword && !form.currentPassword) {
            errors.currentPassword = 'Current password is required';
        }

        if (touched.newPassword) {
            if (!form.newPassword) {
                errors.newPassword = 'New password is required';
            } else if (form.newPassword.length < 8) {
                errors.newPassword = 'Password must be at least 8 characters';
            } else if (!/[A-Z]/.test(form.newPassword)) {
                errors.newPassword = 'Password must contain at least one uppercase letter';
            } else if (!/[a-z]/.test(form.newPassword)) {
                errors.newPassword = 'Password must contain at least one lowercase letter';
            } else if (!/[0-9]/.test(form.newPassword)) {
                errors.newPassword = 'Password must contain at least one number';
            }
        }

        if (touched.confirmPassword) {
            if (!form.confirmPassword) {
                errors.confirmPassword = 'Please confirm your new password';
            } else if (form.newPassword !== form.confirmPassword) {
                errors.confirmPassword = 'Passwords do not match';
            }
        }

        return errors;
    };

    const errors = validateForm();
    const isFormValid =
        form.currentPassword &&
        form.newPassword &&
        form.confirmPassword &&
        Object.keys(errors).length === 0;

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setTouched({
            currentPassword: true,
            newPassword: true,
            confirmPassword: true,
        });

        if (!isFormValid) return;

        setLoading(true);
        setError(null);

        try {
            await changePassword({
                current_password: form.currentPassword,
                new_password: form.newPassword,
            });
            setSuccess(true);
            setTimeout(() => {
                navigate('/profile');
            }, 2000);
        } catch (e: unknown) {
            if (isAxiosError(e)) {
                setError(e.response?.data?.error || 'Failed to change password');
            } else if (e instanceof Error) {
                setError(e.message);
            } else {
                setError('Failed to change password');
            }
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div className="min-h-screen bg-gray-50">
                <Navbar />
                <main className="max-w-md mx-auto px-4 sm:px-6 lg:px-8 py-12">
                    <div className="bg-white rounded-2xl shadow-md border border-gray-100 p-8">
                        <div className="text-center">
                            <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
                                <svg className="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <h3 className="text-lg font-semibold text-gray-900 mb-2">Password changed successfully!</h3>
                            <p className="text-sm text-gray-600">Redirecting to profile...</p>
                        </div>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50">
            <Navbar />
            <main className="max-w-md mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <div className="bg-white rounded-2xl shadow-md border border-gray-100 p-8">
                    <div className="mb-6">
                        <h1 className="text-2xl font-semibold text-gray-900">Change Password</h1>
                        <p className="text-sm text-gray-600 mt-1">Enter your current password and choose a new one</p>
                    </div>

                    {error && (
                        <div className="mb-4 rounded-lg bg-red-50 p-4 text-sm text-red-700 border border-red-200">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-5">
                        <div>
                            <InputField
                                label="Current Password"
                                type={showCurrentPassword ? 'text' : 'password'}
                                name="currentPassword"
                                value={form.currentPassword}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                required
                                rightIcon={
                                    !showCurrentPassword ? (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                                        </svg>
                                    ) : (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                        </svg>
                                    )
                                }
                                onRightIconMouseDown={() => setShowCurrentPassword(true)}
                                onRightIconMouseUp={() => setShowCurrentPassword(false)}
                                onRightIconMouseLeave={() => setShowCurrentPassword(false)}
                            />
                            <FormError message={errors.currentPassword} />
                        </div>

                        <div>
                            <InputField
                                label="New Password"
                                type={showNewPassword ? 'text' : 'password'}
                                name="newPassword"
                                value={form.newPassword}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                required
                                rightIcon={
                                    !showNewPassword ? (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                                        </svg>
                                    ) : (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                        </svg>
                                    )
                                }
                                onRightIconMouseDown={() => setShowNewPassword(true)}
                                onRightIconMouseUp={() => setShowNewPassword(false)}
                                onRightIconMouseLeave={() => setShowNewPassword(false)}
                            />
                            <FormError message={errors.newPassword} />
                            <p className="text-xs text-gray-500 mt-1">
                                Must be at least 8 characters with uppercase, lowercase, and number
                            </p>
                        </div>

                        <div>
                            <InputField
                                label="Confirm New Password"
                                type={showConfirmPassword ? 'text' : 'password'}
                                name="confirmPassword"
                                value={form.confirmPassword}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                required
                                rightIcon={
                                    !showConfirmPassword ? (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                                        </svg>
                                    ) : (
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                        </svg>
                                    )
                                }
                                onRightIconMouseDown={() => setShowConfirmPassword(true)}
                                onRightIconMouseUp={() => setShowConfirmPassword(false)}
                                onRightIconMouseLeave={() => setShowConfirmPassword(false)}
                            />
                            <FormError message={errors.confirmPassword} />
                        </div>

                        <div className="flex gap-3 pt-4">
                            <button
                                type="button"
                                onClick={() => navigate('/profile')}
                                className="flex-1 px-6 py-3 border border-gray-300 text-sm font-medium rounded-xl text-gray-700 bg-white hover:bg-gray-50 transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={loading || !isFormValid}
                                className={`flex-1 px-6 py-3 text-sm font-medium rounded-xl text-white transition-all ${loading || !isFormValid
                                        ? 'bg-gray-300 cursor-not-allowed'
                                        : 'bg-gradient-to-r from-blue-600 to-blue-700 hover:shadow-lg hover:from-blue-700 hover:to-blue-800 active:scale-95'
                                    }`}
                            >
                                {loading ? (
                                    <span className="flex items-center justify-center">
                                        <LoadingSpinner />
                                        <span className="ml-2">Changing...</span>
                                    </span>
                                ) : (
                                    'Change Password'
                                )}
                            </button>
                        </div>
                    </form>
                </div>
            </main>
        </div>
    );
};

export default ChangePassword;
