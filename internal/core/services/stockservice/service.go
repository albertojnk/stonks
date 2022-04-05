package stockservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/albertojnk/stonks/internal/common"
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
	"github.com/albertojnk/stonks/internal/core/ports"
)

type Service struct {
	stockService ports.StockRepository
}

func New(repository ports.StockRepository) *Service {
	return &Service{
		stockService: repository,
	}
}

func (s *Service) Get(ctx *context.Context, stock string, region string) (context.Result, domains.MarketQuotesResponse) {

	stockData := domains.MarketQuotesResponse{}

	stocks := strings.Split(stock, ",")

	_, symbols := s.GetSymbols(ctx, stocks, region)

	url := fmt.Sprintf("https://%v/market/v2/get-quotes?symbols=%v&region=%v", common.GetEnv("rapidapi_host", ""), strings.Join(symbols, ","), region)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Host", common.GetEnv("rapidapi_host", ""))
	req.Header.Add("X-RapidAPI-Key", common.GetEnv("rapidapi_key", ""))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.Logger.Errorf("error doing request, req: %v, err: %v", req, err)
		return ctx.ResultError(res.StatusCode, err.Error(), err), stockData
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &stockData)
	if err != nil {
		ctx.Logger.Errorf("error while unmarshaling autocomplete, err: %v", err)
		return ctx.ResultError(400, "error while unmarshaling autocomplete", err), stockData
	}

	return ctx.ResultSuccess(), stockData
}

func (s *Service) GetSymbols(ctx *context.Context, stocks []string, region string) (context.Result, []string) {

	symbols := []string{}

	for _, stock := range stocks {

		url := fmt.Sprintf("https://%v/auto-complete?q=%v&region=%v", common.GetEnv("rapidapi_host", ""), stock, region)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Host", common.GetEnv("rapidapi_host", ""))
		req.Header.Add("X-RapidAPI-Key", common.GetEnv("rapidapi_key", ""))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			ctx.Logger.Errorf("error doing request, req: %v, err: %v", req, err)
			return ctx.ResultError(res.StatusCode, err.Error(), err), []string{""}
		}

		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		data := domains.AutoCompleteResponse{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			ctx.Logger.Errorf("error while unmarshaling autocomplete, err: %v", err)
			return ctx.ResultError(400, "error while unmarshaling autocomplete", err), []string{""}
		}

		if len(data.Quotes) > 0 {
			symbols = append(symbols, data.Quotes[0].Symbol)
		}
	}

	return ctx.ResultSuccess(), symbols
}
