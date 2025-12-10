import { useState } from 'react';
import { useTranslation } from 'react-i18next';

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
  const { t } = useTranslation('common');

  const validateAll = (address: AddressForm): boolean => {
    const newErrors: CheckoutErrors = {};

    // Address
    if (!address.fullName.trim()) newErrors.fullName = t('checkout.validation.fullNameRequired');
    if (!address.address1.trim()) newErrors.address1 = t('checkout.validation.addressRequired');
    if (!address.number.trim()) newErrors.number = t('checkout.validation.numberRequired');
    if (!address.city.trim()) newErrors.city = t('checkout.validation.cityRequired');
    if (!address.state.trim()) newErrors.state = t('checkout.validation.stateRequired');
    if (!address.postalCode.trim()) newErrors.postalCode = t('checkout.validation.postalCodeRequired');
    if (!address.country.trim()) newErrors.country = t('checkout.validation.countryRequired');

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const clearErrors = () => setErrors({});

  return { errors, setErrors, validateAll, clearErrors };
};
