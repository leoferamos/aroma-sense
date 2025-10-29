import React, { useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import InputField from '../components/InputField';
import FormError from '../components/FormError';
import CartItem from '../components/CartItem';
import { useCart } from '../contexts/CartContext';
import { formatCurrency } from '../utils/format';
import { useCheckoutValidation, type AddressForm, type PaymentForm } from '../hooks/useCheckoutValidation';

const Checkout: React.FC = () => {
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


  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    if (cartIsEmpty) return;
  if (!validateAll(address, payment)) return;
    setSubmitting(true);
    // Simulate success and redirect to confirmation screen
    setTimeout(() => {
      navigate('/order-confirmation', {
        replace: true,
        state: {
          orderTotal: cart?.total ?? 0,
          itemsCount: cart?.item_count ?? 0,
          customerName: address.fullName,
        },
      });
    }, 600);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>

        {loading ? (
          <div className="text-gray-600">Loading your cart...</div>
        ) : cartIsEmpty ? (
          <div className="bg-white shadow rounded-lg p-8 text-center">
            <p className="text-gray-700 mb-4">Your cart is empty.</p>
            <Link to="/products" className="text-blue-600 hover:underline font-medium">Continue shopping</Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Left: Forms */}
            <form onSubmit={handleSubmit} noValidate className="lg:col-span-2 space-y-8">
              {/* Address */}
              <section className="bg-white shadow rounded-lg p-6">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Shipping address</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="sm:col-span-2">
                    <InputField
                      label="Full name"
                      name="fullName"
                      value={address.fullName}
                      onChange={(e) => setAddress({ ...address, fullName: e.target.value })}
                      onBlur={() => setErrors((prev) => ({ ...prev, fullName: prev.fullName }))}
                    />
                    <FormError message={errors.fullName} />
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
                      label="Postal code"
                      name="postalCode"
                      value={address.postalCode}
                      onChange={(e) => setAddress({ ...address, postalCode: e.target.value })}
                      autoComplete="postal-code"
                    />
                    <FormError message={errors.postalCode} />
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

              {/* Payment */}
              <section className="bg-white shadow rounded-lg p-6">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Payment</h2>
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
                      onChange={(e) => setPayment({ ...payment, cardNumber: e.target.value })}
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
                      onChange={(e) => setPayment({ ...payment, expiry: e.target.value })}
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
                      onChange={(e) => setPayment({ ...payment, cvc: e.target.value })}
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
                  disabled={submitting || cartIsEmpty}
                  aria-disabled={submitting || cartIsEmpty}
                  aria-busy={submitting}
                  className={`inline-flex items-center justify-center rounded-md bg-blue-600 px-6 py-3 text-white font-medium shadow hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed`}
                >
                  {submitting ? 'Placing order...' : 'Place Order'}
                </button>
              </div>
            </form>

            {/* Right: Cart Summary */}
            <aside className="lg:col-span-1">
              <div className="bg-white shadow rounded-lg p-6 sticky top-24">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Order summary</h2>
                
                {/* Error display */}
                {error && (
                  <div className="mb-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700" role="alert">
                    {error}
                  </div>
                )}
                
                <ul className="divide-y divide-gray-200 mb-4">
                  {cart!.items.map((item) => (
                    <CartItem
                      key={item.id}
                      item={item}
                      onRemove={removeItem}
                      isRemoving={isRemovingItem(item.id)}
                    />
                  ))}
                </ul>
                <div className="flex justify-between text-gray-700">
                  <span>Subtotal</span>
                  <span className="font-semibold">{formatCurrency(cart!.total)}</span>
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
