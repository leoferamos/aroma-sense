package shipping

import "github.com/leoferamos/aroma-sense/internal/dto"

// mapProviderQuotes converts provider quotes into public DTO options, filtering invalid entries.
func mapProviderQuotes(items []providerQuote) []dto.ShippingOption {
	if len(items) == 0 {
		return []dto.ShippingOption{}
	}
	opts := make([]dto.ShippingOption, 0, len(items))
	for _, it := range items {
		if it.HasError || it.Price <= 0 {
			continue
		}
		carrier := it.Company.Name
		serviceCode := it.Name
		opts = append(opts, dto.ShippingOption{
			Carrier:       carrier,
			ServiceCode:   serviceCode,
			Price:         it.Price,
			EstimatedDays: it.DeliveryTime,
		})
	}
	return opts
}
