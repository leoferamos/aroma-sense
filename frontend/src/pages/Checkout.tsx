import React, { useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import BackButton from '../components/BackButton';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import CartItem from '../components/CartItem';
import { useCart } from '../hooks/useCart';
import { formatCurrency } from '../utils/format';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorState from '../components/ErrorState';
import { useCheckoutValidation, type AddressForm, type PaymentForm } from '../hooks/useCheckoutValidation';
import useShippingOptions from '../hooks/useShippingOptions';
import useCepLookup from '../hooks/useCepLookup';
import type { CartItem as CartItemType } from '../types/cart';
import { useTranslation } from 'react-i18next';
import type { ShippingOption } from '../types/shipping';
import { createOrder, type OrderCreateRequest } from '../services/order';

const Checkout: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { cart, loading, removeItem, error, isRemovingItem } = useCart();

  const [address, setAddress] = useState<AddressForm>({
    fullName: '',
    address1: '',
    address2: '',
    city: '',
    state: '',
    postalCode: '',
    country: 'Brazil',
  });

  const [payment, setPayment] = useState<PaymentForm>({
    cardName: '',
    cardNumber: '',
    expiry: '',
    cvc: '',
  });

  const { errors, validateAll, setErrors } = useCheckoutValidation();
  const [submitting, setSubmitting] = useState(false);
  const cartIsEmpty = useMemo(() => !cart || cart.items.length === 0, [cart]);
  const { options: shippingOptions, loading: shippingLoading, error: shippingError } = useShippingOptions(address.postalCode);
  const [selectedShipping, setSelectedShipping] = useState<ShippingOption | null>(null);
  const { lookupCep, loading: cepLoading, error: cepError } = useCepLookup();


  const handleSubmit: React.FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    if (cartIsEmpty) return;
    if (!validateAll(address, payment)) return;
    if (!selectedShipping) {
      setErrors((prev) => ({ ...prev, postalCode: prev.postalCode || 'Select a shipping option' }));
      return;
    }
    setSubmitting(true);
    try {
      const shipping_address = `${address.address1}${address.address2 ? ', ' + address.address2 : ''}, ${address.city} - ${address.state}, ${address.postalCode}`;
      const payload: OrderCreateRequest = {
        payment_method: 'pix',
        shipping_address,
        shipping_selection: {
          carrier: selectedShipping.carrier,
          service_code: selectedShipping.service_code,
          price: selectedShipping.price,
          estimated_days: selectedShipping.estimated_days,
          quote_id: null,
        },
      };
      const order = await createOrder(payload);
      navigate('/order-confirmation', {
        replace: true,
        state: {
          orderTotal: order.total_amount,
          itemsCount: order.item_count,
          customerName: address.fullName,
        },
      });
    } catch (err) {
      console.error('Failed to create order', err);
      setErrors((prev) => ({ ...prev, address1: t('errors.failedToCreateOrder') }));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div className="mb-4">
          <BackButton fallbackPath="/products" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('cart.checkout')}</h1>

        {loading ? (
          <LoadingSpinner message="Loading your cart..." />
        ) : cartIsEmpty ? (
          <div className="bg-white rounded-xl shadow-sm p-8 text-center border border-gray-100">
            <p className="text-gray-700 mb-4">{t('checkout.cartEmpty')}</p>
            <Link to="/products" className="text-blue-600 hover:underline font-medium">{t('checkout.continueShopping')}</Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Left: Forms */}
            <form onSubmit={handleSubmit} noValidate className="lg:col-span-2 space-y-8">
              {/* Address */}
              <section className="bg-white shadow-sm rounded-xl p-6 border border-gray-100">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('checkout.shippingAddress')}</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="sm:col-span-2">
                    <InputField
                      label={t('checkout.fullName')}
                      name="fullName"
                      value={address.fullName}
                      onChange={(e) => setAddress({ ...address, fullName: e.target.value })}
                      onBlur={() => setErrors((prev) => ({ ...prev, fullName: prev.fullName }))}
                    />
                    <FormError message={errors.fullName} />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <InputField
                        label="Postal code"
                        name="postalCode"
                        value={address.postalCode}
                        onChange={(e) => {
                          const cleaned = e.target.value.replace(/\D/g, '').slice(0, 8);
                          setAddress({ ...address, postalCode: cleaned });
                          if (cleaned.length === 8) {
                            void (async () => {
                              const data = await lookupCep(cleaned);
                              if (data) {
                                setAddress((prev) => ({
                                  ...prev,
                                  address1: data.logradouro || prev.address1,
                                  city: data.localidade || prev.city,
                                  state: data.uf || prev.state,
                                  postalCode: cleaned,
                                }));
                              }
                            })();
                          }
                        }}
                        onBlur={() => {
                          void (async () => {
                            const data = await lookupCep(address.postalCode);
                            if (data) {
                              setAddress((prev) => ({
                                ...prev,
                                address1: data.logradouro || prev.address1,
                                city: data.localidade || prev.city,
                                state: data.uf || prev.state,
                                postalCode: (address.postalCode || '').replace(/\D/g, ''),
                              }));
                            }
                          })();
                        }}
                        autoComplete="postal-code"
                      />
                      {cepLoading && <span className="text-sm text-gray-500">Buscando...</span>}
                    </div>
                    <FormError message={cepError || errors.postalCode} />
                  </div>
                  <div className="sm:col-span-2">
                    <InputField
                      label="Address line 1"
                      name="address1"
                      value={address.address1}
                      onChange={(e) => setAddress({ ...address, address1: e.target.value })}
                      autoComplete="address-line1"
                    />
                    <FormError message={errors.address1} />
                  </div>
                  <div className="sm:col-span-2">
                    <InputField
                      label="Address line 2 (optional)"
                      name="address2"
                      value={address.address2 || ''}
                      onChange={(e) => setAddress({ ...address, address2: e.target.value })}
                      autoComplete="address-line2"
                    />
                  </div>
                  <div>
                    <InputField
                      label="City"
                      name="city"
                      value={address.city}
                      onChange={(e) => setAddress({ ...address, city: e.target.value })}
                      autoComplete="address-level2"
                    />
                    <FormError message={errors.city} />
                  </div>
                  <div>
                    <InputField
                      label="State/Region"
                      name="state"
                      value={address.state}
                      onChange={(e) => setAddress({ ...address, state: e.target.value })}
                      autoComplete="address-level1"
                    />
                    <FormError message={errors.state} />
                  </div>
                  <div>
                    <InputField
                      label="Country"
                      name="country"
                      value={address.country}
                      onChange={(e) => setAddress({ ...address, country: e.target.value })}
                      autoComplete="country-name"
                    />
                    <FormError message={errors.country} />
                  </div>
                </div>
              </section>

              {/* Shipping options */}
              <section className="bg-white shadow-sm rounded-xl p-6 border border-gray-100">
                <h2 className="text-xl font-semibold text-gray-900 mb-2">Shipping</h2>
                <p className="text-sm text-gray-600 mb-4">Enter your postal code to see available options.</p>
                {shippingError && <ErrorState message={shippingError} />}
                {shippingLoading && <LoadingSpinner message="Loading shipping options..." />}
                {!shippingLoading && shippingOptions.length > 0 && (
                  <div className="space-y-2">
                    {shippingOptions.map((opt, idx) => (
                      <label key={`${opt.carrier}-${opt.service_code}-${idx}`} className={`flex items-center justify-between border rounded-md p-3 cursor-pointer ${selectedShipping === opt ? 'border-blue-500 ring-1 ring-blue-300' : 'border-gray-200'}`}>
                        <div className="flex items-center gap-3">
                          <input
                            type="radio"
                            name="shippingOption"
                            value={`${opt.carrier}:${opt.service_code}`}
                            checked={selectedShipping?.carrier === opt.carrier && selectedShipping?.service_code === opt.service_code}
                            onChange={() => setSelectedShipping(opt)}
                          />
                          <div>
                            <div className="font-medium text-gray-900">{opt.carrier} — {opt.service_code}</div>
                            <div className="text-sm text-gray-600">ETA: {opt.estimated_days} day{opt.estimated_days === 1 ? '' : 's'}</div>
                          </div>
                        </div>
                        <div className="font-semibold">{formatCurrency(opt.price)}</div>
                      </label>
                    ))}
                  </div>
                )}
                {!shippingLoading && shippingOptions.length === 0 && address.postalCode && address.postalCode.replace(/\D/g, '').length >= 8 && (
                  <p className="text-sm text-gray-600">No shipping options for this postal code.</p>
                )}
              </section>

              {/* Payment */}
              <section className="bg-white shadow-sm rounded-xl p-6 border border-gray-100">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('checkout.payment')}</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="sm:col-span-2">
                    <InputField
                      label="Name on card"
                      name="cardName"
                      value={payment.cardName}
                      onChange={(e) => setPayment({ ...payment, cardName: e.target.value })}
                      autoComplete="cc-name"
                    />
                    <FormError message={errors.cardName} />
                  </div>
                  <div className="sm:col-span-2">
                    <InputField
                      label="Card number"
                      name="cardNumber"
                      value={payment.cardNumber}
                      onChange={(e) => {
                        const digitsOnly = e.target.value.replace(/\D/g, '').slice(0, 19);
                        const formatted = digitsOnly.replace(/(\d{4})(?=\d)/g, '$1 ');
                        setPayment({ ...payment, cardNumber: formatted });
                      }}
                      placeholder="1234 5678 9012 3456"
                      autoComplete="cc-number"
                    />
                    <FormError message={errors.cardNumber} />
                  </div>
                  <div>
                    <InputField
                      label="Expiry (MM/YY)"
                      name="expiry"
                      value={payment.expiry}
                      onChange={(e) => {
                        let value = e.target.value.replace(/\D/g, '').slice(0, 4);
                        if (value.length >= 2) {
                          value = value.slice(0, 2) + '/' + value.slice(2, 4);
                        }
                        setPayment({ ...payment, expiry: value });
                      }}
                      placeholder="MM/YY"
                      autoComplete="cc-exp"
                    />
                    <FormError message={errors.expiry} />
                  </div>
                  <div>
                    <InputField
                      label="CVC"
                      name="cvc"
                      value={payment.cvc}
                      onChange={(e) => {
                        const value = e.target.value.replace(/\D/g, '').slice(0, 4);
                        setPayment({ ...payment, cvc: value });
                      }}
                      placeholder="123"
                      autoComplete="cc-csc"
                    />
                    <FormError message={errors.cvc} />
                  </div>
                </div>
              </section>

              <div className="flex justify-end">
                <button
                  type="submit"
                  disabled={submitting || cartIsEmpty || !selectedShipping}
                  aria-disabled={submitting || cartIsEmpty || !selectedShipping}
                  aria-busy={submitting}
                  className={`inline-flex items-center justify-center rounded-lg bg-blue-600 px-6 py-3 text-white font-semibold shadow-sm hover:bg-blue-700 hover:shadow-md transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed`}
                >
                  {submitting ? t('checkout.placingOrder') : t('checkout.placeOrder')}
                </button>
              </div>
            </form>

            {/* Right: Cart Summary */}
            <aside className="lg:col-span-1">
              <div className="bg-white shadow-sm rounded-xl p-6 sticky top-24 border border-gray-100">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('checkout.orderSummary')}</h2>

                {/* Error display */}
                {error && (
                  <ErrorState message={error} />
                )}

                <ul className="divide-y divide-gray-200 mb-4">
                  {cart!.items.map((item: CartItemType) => (
                    <CartItem
                      key={item.id}
                      item={item}
                      onRemove={removeItem}
                      isRemoving={isRemovingItem(item.id)}
                      showQuantityControls={true}
                    />
                  ))}
                </ul>
                <div className="flex justify-between text-gray-700">
                  <span>{t('cart.subtotal')}</span>
                  <span className="font-semibold">{formatCurrency(cart!.total)}</span>
                </div>
                <div className="flex justify-between text-gray-700 mt-2">
                  <span>Shipping</span>
                  <span className="font-semibold">{selectedShipping ? formatCurrency(selectedShipping.price) : '-'}</span>
                </div>
                <div className="flex justify-between text-gray-900 mt-2 border-t pt-2">
                  <span>{t('cart.total')}</span>
                  <span className="font-bold">{formatCurrency(cart!.total + (selectedShipping?.price ?? 0))}</span>
                </div>
              </div>
            </aside>
          </div>
        )}
      </main>
    </div>
  );
};

export default Checkout;
