package handlers

import (
	"github.com/albertojnk/stonks/internal/context"
	"github.com/gin-gonic/gin"
)

func (hdl *HTTPHandler) StockGet(ctx *context.Context, c *gin.Context) error {

	stock := c.Query("stock")
	if stock == "" {
		ctx.ResultError(1, "Invalid Request").JSON(c, nil)
		return nil
	}

	region := c.Query("region")
	if region == "" {
		region = "BR"
	}

	res, body := hdl.stockService.Get(ctx, stock, region)

	res.JSON(c, body)

	return nil
}
