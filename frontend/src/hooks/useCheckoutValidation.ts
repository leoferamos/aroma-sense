import { useState } from 'react';
import { messages } from '../constants/messages';

export interface AddressForm {
  fullName: string;
  address1: string;
  address2?: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
}

export interface PaymentForm {
  cardName: string;
  cardNumber: string;
  expiry: string; 
  cvc: string;
}

export type CheckoutErrors = Partial<
  Record<
    | keyof AddressForm
    | keyof PaymentForm,
    string
  >
>;

export const useCheckoutValidation = () => {
  const [errors, setErrors] = useState<CheckoutErrors>({});

  const validateAll = (address: AddressForm, payment: PaymentForm): boolean => {
    const newErrors: CheckoutErrors = {};

    // Address
    if (!address.fullName.trim()) newErrors.fullName = messages.fullNameRequired;
    if (!address.address1.trim()) newErrors.address1 = messages.addressRequired;
    if (!address.city.trim()) newErrors.city = messages.cityRequired;
    if (!address.state.trim()) newErrors.state = messages.stateRequired;
    if (!address.postalCode.trim()) newErrors.postalCode = messages.postalCodeRequired;
    if (!address.country.trim()) newErrors.country = messages.countryRequired;

    // Payment
    if (!payment.cardName.trim()) newErrors.cardName = messages.nameOnCardRequired;
    const digitsOnly = payment.cardNumber.replace(/\s+/g, '');
    if (!/^\d{13,19}$/.test(digitsOnly)) newErrors.cardNumber = messages.cardNumberInvalid;
    if (!/^(0[1-9]|1[0-2])\/(\d{2})$/.test(payment.expiry)) newErrors.expiry = messages.expiryInvalid;
    if (!/^\d{3,4}$/.test(payment.cvc)) newErrors.cvc = messages.cvcInvalid;

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const clearErrors = () => setErrors({});

  return { errors, setErrors, validateAll, clearErrors };
};
