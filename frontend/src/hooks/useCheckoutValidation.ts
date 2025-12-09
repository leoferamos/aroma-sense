import { useState } from 'react';
import { messages } from '../constants/messages';

export interface AddressForm {
  fullName: string;
  address1: string;
  number: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
}

export type CheckoutErrors = Partial<
  Record<
    keyof AddressForm,
    string
  >
>;

export const useCheckoutValidation = () => {
  const [errors, setErrors] = useState<CheckoutErrors>({});

  const validateAll = (address: AddressForm): boolean => {
    const newErrors: CheckoutErrors = {};

    // Address
    if (!address.fullName.trim()) newErrors.fullName = messages.fullNameRequired;
    if (!address.address1.trim()) newErrors.address1 = messages.addressRequired;
    if (!address.number.trim()) newErrors.number = messages.numberRequired;
    if (!address.city.trim()) newErrors.city = messages.cityRequired;
    if (!address.state.trim()) newErrors.state = messages.stateRequired;
    if (!address.postalCode.trim()) newErrors.postalCode = messages.postalCodeRequired;
    if (!address.country.trim()) newErrors.country = messages.countryRequired;

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const clearErrors = () => setErrors({});

  return { errors, setErrors, validateAll, clearErrors };
};
