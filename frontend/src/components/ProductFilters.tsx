import React from 'react';
import { useTranslation } from 'react-i18next';
import type { Dispatch, SetStateAction } from 'react';

interface FilterState {
    genders: string[];
    priceRange: [number, number];
}

interface ProductFiltersProps {
    filters: FilterState;
    onFiltersChange: Dispatch<SetStateAction<FilterState>>;
    minPrice?: number;
    maxPrice?: number;
}const ProductFilters: React.FC<ProductFiltersProps> = ({
    filters,
    onFiltersChange,
    minPrice = 0,
    maxPrice = 1000,
}) => {
    const { t } = useTranslation('common');

    const genderOptions = [
        { value: 'masculino', label: t('filters.masculine') || 'Masculino' },
        { value: 'feminino', label: t('filters.feminine') || 'Feminino' },
        { value: 'unissex', label: t('filters.unisex') || 'Unissex' },
    ];

    const handleGenderChange = (gender: string) => {
        const newGenders = filters.genders.includes(gender)
            ? filters.genders.filter((g) => g !== gender)
            : [...filters.genders, gender];
        onFiltersChange({ ...filters, genders: newGenders });
    };

    const handleMinPriceChange = (value: number) => {
        const newRange: [number, number] = [value, filters.priceRange[1]];
        onFiltersChange({ ...filters, priceRange: newRange });
    };

    const handleMaxPriceChange = (value: number) => {
        const newRange: [number, number] = [filters.priceRange[0], value];
        onFiltersChange({ ...filters, priceRange: newRange });
    };

    const handleResetFilters = () => {
        onFiltersChange({
            genders: [],
            priceRange: [minPrice, maxPrice],
        });
    };

    const hasActiveFilters =
        filters.genders.length > 0 ||
        filters.priceRange[0] !== minPrice ||
        filters.priceRange[1] !== maxPrice;

    return (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-semibold text-gray-900">
                    {t('filters.title') || 'Filtros'}
                </h2>
                {hasActiveFilters && (
                    <button
                        onClick={handleResetFilters}
                        className="text-sm text-purple-600 hover:text-purple-700 font-medium"
                    >
                        {t('filters.reset') || 'Limpar'}
                    </button>
                )}
            </div>

            {/* Gender Filter */}
            <div className="mb-8">
                <h3 className="text-sm font-semibold text-gray-900 mb-3">
                    {t('filters.gender') || 'Gênero'}
                </h3>
                <div className="space-y-2">
                    {genderOptions.map((option) => (
                        <label key={option.value} className="flex items-center cursor-pointer">
                            <input
                                type="checkbox"
                                checked={filters.genders.includes(option.value)}
                                onChange={() => handleGenderChange(option.value)}
                                className="w-4 h-4 text-purple-600 rounded border-gray-300 focus:ring-purple-500"
                            />
                            <span className="ml-3 text-sm text-gray-600 hover:text-gray-900">
                                {option.label}
                            </span>
                        </label>
                    ))}
                </div>
            </div>

            {/* Price Range Filter */}
            <div className="mb-8">
                <h3 className="text-sm font-semibold text-gray-900 mb-4">
                    {t('filters.price') || 'Preço'}
                </h3>
                <div className="space-y-4">
                    {/* Min Price */}
                    <div>
                        <label className="block text-xs text-gray-600 mb-2">
                            {t('filters.minPrice') || 'Preço Mínimo'}
                        </label>
                        <input
                            type="range"
                            min={minPrice}
                            max={maxPrice}
                            value={filters.priceRange[0]}
                            onChange={(e) => handleMinPriceChange(Number(e.target.value))}
                            className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
                        />
                        <div className="mt-1 text-sm font-medium text-gray-900">
                            R$ {filters.priceRange[0].toFixed(2)}
                        </div>
                    </div>

                    {/* Max Price */}
                    <div>
                        <label className="block text-xs text-gray-600 mb-2">
                            {t('filters.maxPrice') || 'Preço Máximo'}
                        </label>
                        <input
                            type="range"
                            min={minPrice}
                            max={maxPrice}
                            value={filters.priceRange[1]}
                            onChange={(e) => handleMaxPriceChange(Number(e.target.value))}
                            className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-purple-600"
                        />
                        <div className="mt-1 text-sm font-medium text-gray-900">
                            R$ {filters.priceRange[1].toFixed(2)}
                        </div>
                    </div>

                    {/* Price Range Display */}
                    <div className="bg-gray-50 rounded p-3 mt-4">
                        <p className="text-sm text-gray-600">
                            {t('filters.range') || 'Intervalo de preço'}:{' '}
                            <span className="font-semibold text-gray-900">
                                R$ {filters.priceRange[0].toFixed(2)} - R$ {filters.priceRange[1].toFixed(2)}
                            </span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProductFilters;
