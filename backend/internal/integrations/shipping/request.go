package shipping

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/validation"
)

// buildQuoteRequest creates the provider request payload from domain inputs.
func (p *Provider) buildQuoteRequest(originPostalCode, destPostalCode string, parcels []model.Parcel, insuredValue float64) quoteRequest {
	var reqBody quoteRequest
	reqBody.From.PostalCode = validation.NormalizeCEP(originPostalCode)
	reqBody.To.PostalCode = validation.NormalizeCEP(destPostalCode)
	if len(parcels) > 0 {
		pr := parcels[0]
		reqBody.Package.Weight = pr.WeightKg
		reqBody.Package.Height = pr.HeightCm
		reqBody.Package.Width = pr.WidthCm
		reqBody.Package.Length = pr.LengthCm
	}
	reqBody.Services = p.services
	reqBody.Options.OwnHand = false
	reqBody.Options.Receipt = false
	reqBody.Options.InsuranceValue = insuredValue
	reqBody.Options.UseInsuranceValue = insuredValue > 0
	return reqBody
}
