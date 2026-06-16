package service

import (
	"strings"

	"github.com/GMWalletApp/epusdt/config"
	"github.com/GMWalletApp/epusdt/model/data"
	"github.com/GMWalletApp/epusdt/model/mdb"
	"github.com/GMWalletApp/epusdt/model/response"
	"github.com/GMWalletApp/epusdt/util/constant"
)

// ErrOrder is returned when checkout initialization cannot find the trade id.
var ErrOrder = constant.OrderNotExists

// GetCheckoutCounterByTradeId returns checkout initialization data for an existing order.
// It does not decide the payment state; callers should use CheckStatus for that.
func GetCheckoutCounterByTradeId(tradeId string) (*response.CheckoutCounterResponse, error) {
	orderInfo, err := data.GetOrderInfoByTradeId(tradeId)
	if err != nil {
		return nil, err
	}
	if orderInfo.ID <= 0 {
		return nil, ErrOrder
	}
	paymentType := mdb.PaymentTypeGmpay
	if strings.EqualFold(orderInfo.PaymentType, mdb.PaymentTypeEpay) {
		paymentType = mdb.PaymentTypeEpay
	}

	resp := &response.CheckoutCounterResponse{
		TradeId:        orderInfo.TradeId,
		Amount:         orderInfo.Amount,
		ActualAmount:   orderInfo.ActualAmount,
		Token:          orderInfo.Token,
		Currency:       orderInfo.Currency,
		ReceiveAddress: orderInfo.ReceiveAddress,
		Network:        orderInfo.Network,
		Status:         orderInfo.Status,
		PaymentType:    strings.ToLower(paymentType),
		ExpirationTime: orderInfo.CreatedAt.AddMinutes(config.GetOrderExpirationTime()).TimestampMilli(),
		RedirectUrl:    orderInfo.RedirectUrl,
		CreatedAt:      orderInfo.CreatedAt.TimestampMilli(),
		IsSelected:     orderInfo.IsSelected,
	}
	if orderInfo.PayProvider == mdb.PaymentProviderOkPay {
		providerRow, rowErr := data.GetProviderOrderByTradeIDAndProvider(orderInfo.TradeId, mdb.PaymentProviderOkPay)
		if rowErr != nil {
			return nil, rowErr
		}
		if providerRow.ID > 0 {
			resp.PaymentUrl = providerRow.PayURL
		}
	}
	return resp, nil
}
