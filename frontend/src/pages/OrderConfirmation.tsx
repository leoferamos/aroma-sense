import React from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import Navbar from '../components/Navbar';
import BackButton from '../components/BackButton';
import { formatCurrency } from '../utils/format';

type LocationState = {
  orderTotal?: number;
  itemsCount?: number;
  customerName?: string;
};

const OrderConfirmation: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation('common');
  const state = (location.state || {}) as LocationState;

  const total = state.orderTotal ?? 0;
  const items = state.itemsCount ?? 0;
  const name = state.customerName ?? '';

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div className="mb-4">
          <BackButton fallbackPath="/products" />
        </div>
        <div className="bg-white shadow rounded-lg p-8 text-center">
          <div className="mx-auto h-16 w-16 rounded-full bg-green-100 flex items-center justify-center mb-4">
            <svg className="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900">{t('checkout.orderConfirmation.thankYou')}{name ? `, ${name}` : ''}!</h1>
          <p className="mt-2 text-gray-700">{t('checkout.orderConfirmation.orderPlaced')}</p>
          <div className="mt-6 inline-flex items-center gap-6 text-gray-800">
            <div>
              <div className="text-sm text-gray-500">{t('checkout.orderConfirmation.items')}</div>
              <div className="text-lg font-semibold">{items}</div>
            </div>
            <div className="h-8 w-px bg-gray-200" />
            <div>
              <div className="text-sm text-gray-500">{t('checkout.orderConfirmation.total')}</div>
              <div className="text-lg font-semibold">{formatCurrency(total)}</div>
            </div>
          </div>
          <div className="mt-8 flex items-center justify-center gap-4">
            <Link to="/products" className="px-6 py-3 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700 transition-colors">{t('checkout.orderConfirmation.continueShopping')}</Link>
            <button onClick={() => navigate(-1)} className="px-6 py-3 rounded-md border border-gray-300 text-gray-700 font-medium hover:bg-gray-50 transition-colors">{t('checkout.orderConfirmation.goBack')}</button>
          </div>
        </div>
      </main>
    </div>
  );
};

export default OrderConfirmation;
