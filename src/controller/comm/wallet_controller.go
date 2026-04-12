package comm

import (
	"strconv"

	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/mdb"
	"github.com/assimon/luuu/util/constant"
	"github.com/labstack/echo/v4"
)

type addWalletRequest struct {
	Network string `json:"network" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type changeStatusRequest struct {
	Status int `json:"status" validate:"required|in:1,2"`
}

// AddWallet 添加钱包地址
func (c *BaseCommController) AddWallet(ctx echo.Context) error {
	req := new(addWalletRequest)
	if err := ctx.Bind(req); err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	if err := c.ValidateStruct(ctx, req); err != nil {
		return c.FailJson(ctx, err)
	}
	wallet, err := data.AddWalletAddressWithNetwork(req.Network, req.Address)
	if err != nil {
		return c.FailJson(ctx, err)
	}
	return c.SucJson(ctx, wallet)
}

// ListWallets 获取钱包列表
func (c *BaseCommController) ListWallets(ctx echo.Context) error {
	network := ctx.QueryParam("network")
	var wallets []mdb.WalletAddress
	var err error
	if network != "" {
		wallets, err = data.GetAllWalletAddressByNetwork(network)
	} else {
		wallets, err = data.GetAllWalletAddress()
	}
	if err != nil {
		return c.FailJson(ctx, err)
	}
	return c.SucJson(ctx, wallets)
}

// GetWallet 获取单个钱包
func (c *BaseCommController) GetWallet(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	wallet, err := data.GetWalletAddressById(id)
	if err != nil {
		return c.FailJson(ctx, err)
	}
	if wallet.ID <= 0 {
		return c.FailJson(ctx, constant.NotAvailableWalletAddress)
	}
	return c.SucJson(ctx, wallet)
}

// ChangeWalletStatus 启用/禁用钱包
func (c *BaseCommController) ChangeWalletStatus(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	req := new(changeStatusRequest)
	if err := ctx.Bind(req); err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	if err := c.ValidateStruct(ctx, req); err != nil {
		return c.FailJson(ctx, err)
	}
	if err := data.ChangeWalletAddressStatus(id, req.Status); err != nil {
		return c.FailJson(ctx, err)
	}
	return c.SucJson(ctx, nil)
}

// DeleteWallet 删除钱包
func (c *BaseCommController) DeleteWallet(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	if err := data.DeleteWalletAddressById(id); err != nil {
		return c.FailJson(ctx, err)
	}
	return c.SucJson(ctx, nil)
}
