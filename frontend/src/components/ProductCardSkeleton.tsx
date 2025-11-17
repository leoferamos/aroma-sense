import React from "react";

const ProductCardSkeleton: React.FC = () => {
  return (
    <div className="bg-white/90 rounded-xl border border-gray-100 shadow-lg overflow-hidden animate-pulse">
      <div className="h-60 bg-gray-100" />
      <div className="p-6">
        <div className="h-3 w-24 bg-gray-200 rounded mb-2" />
        <div className="h-5 w-3/4 bg-gray-200 rounded mb-3" />
        <div className="flex items-center gap-2 mb-4">
          <div className="h-4 w-16 bg-gray-200 rounded-full" />
          <div className="h-4 w-16 bg-gray-200 rounded-full" />
        </div>
        <div className="flex items-center justify-between mb-4">
          <div className="h-6 w-24 bg-gray-200 rounded" />
          <div className="h-3 w-20 bg-gray-200 rounded" />
        </div>
        <div className="h-10 w-full bg-gray-200 rounded" />
      </div>
    </div>
  );
};

export default ProductCardSkeleton;
