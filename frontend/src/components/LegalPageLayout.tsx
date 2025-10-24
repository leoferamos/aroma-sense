import React from 'react';
import { Link } from 'react-router-dom';
import WordGrid from './WordGrid';

interface LegalPageLayoutProps {
  title: string;
  lastUpdate: string;
  children: React.ReactNode;
}

const LegalPageLayout: React.FC<LegalPageLayoutProps> = ({ title, lastUpdate, children }) => {
  return (
    <div className="min-h-screen flex flex-col md:flex-row">
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

      <div className="w-full md:w-1/2 bg-white flex items-center justify-center px-4 py-8 md:px-0 md:py-0 relative">
        <div className="w-full max-w-3xl px-4 md:px-8 py-8 md:py-12 rounded-lg shadow-md">
          <div className="relative mb-6">
            <Link 
              to="/register" 
              className="absolute left-3 top-1/2 -translate-y-1/2 flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
            >
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                fill="none" 
                viewBox="0 0 24 24" 
                strokeWidth={2} 
                stroke="currentColor" 
                className="w-6 h-6"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
              </svg>
              <span className="text-sm font-medium">Back</span>
            </Link>
            <div className="flex flex-col items-center">
              <img src="/logo.png" alt="Logo" className="h-16 md:h-20 mb-4" />
              <h2 className="text-2xl md:text-3xl font-medium text-center" style={{ fontFamily: 'Poppins, sans-serif' }}>
                {title}
              </h2>
              <span className="text-sm text-gray-500 mt-2">Last updated: {lastUpdate}</span>
            </div>
          </div>

          <div className="prose max-w-none text-gray-800 max-h-[60vh] overflow-y-auto px-6 py-4 rounded-xl bg-gradient-to-b from-gray-50 to-white shadow-inner leading-relaxed">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
};

export default LegalPageLayout;
