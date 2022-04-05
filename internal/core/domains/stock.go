package domains

type AutoCompleteResponse struct {
	Quotes []QuotesResponse `json:"quotes"`
}

type QuotesResponse struct {
	Longname string `json:"longname"`
	Symbol   string `json:"symbol"`
	Sector   string `json:"sector"`
}

type Stock struct {
	Price         StockPrice         `json:"price,omitempty"`
	SummaryDetail StockSummaryDetail `json:"summaryDetail,omitempty"`
	Symbol        string             `json:"symbol,omitempty"`
}

type StockPrice struct {
	MarketOpen                PriceMarketTmp `json:"regularMarketOpen,omitempty"`
	MarketVariation           PriceMarketTmp `json:"regularMarketChange,omitempty"`
	MarketPreviousClose       PriceMarketTmp `json:"regularMarketPreviousClose,omitempty"`
	MarketLow                 PriceMarketTmp `json:"regularMarketDayLow,omitempty"`
	MarketHigh                PriceMarketTmp `json:"regularMarketDayHigh,omitempty"`
	MarketPrice               PriceMarketTmp `json:"regularMarketPrice,omitempty"`
	MarketVariationPercentage PriceMarketTmp `json:"regularMarketChangePercent,omitempty"`
	Symbol                    string         `json:"symbol,omitempty"`
	LongName                  string         `json:"longName,omitempty"`
	ShortName                 string         `json:"shortName,omitempty"`
	Currency                  string         `json:"currency,omitempty"`
}

type StockSummaryDetail struct {
	MarketPreviousClose  PriceMarketTmp `json:"previousClose,omitempty"`
	MarketOpen           PriceMarketTmp `json:"open,omitempty"`
	TwoHundredDayAverage PriceMarketTmp `json:"twoHundredDayAverage,omitempty"`
	AnnualDY             PriceMarketTmp `json:"trailingAnnualDividendYield,omitempty"`
	PayoutRatio          PriceMarketTmp `json:"payoutRatio,omitempty"`
	MarketHigh           PriceMarketTmp `json:"dayHigh,omitempty"`
	FiftyDayAverage      PriceMarketTmp `json:"fiftyDayAverage,omitempty"`
	MarketLow            PriceMarketTmp `json:"dayLow,omitempty"`
}

type MarketQuotesResponse struct {
	QuoteResponse MarketQuotes `json:"quoteResponse"`
}

type MarketQuotes struct {
	Result []MarketResult `json:"result"`
}

type MarketResult struct {
	MarketPrice               float64 `json:"regularMarketPrice,omitempty"`
	MarketHigh                float64 `json:"regularMarketDayHigh,omitempty"`
	MarketDayRange            string  `json:"regularMarketDayRange,omitempty"`
	MarketLow                 float64 `json:"regularMarketDayLow,omitempty"`
	MarketPreviousClose       float64 `json:"regularMarketPreviousClose,omitempty"`
	MarketVariation           float64 `json:"regularMarketChange,omitempty"`
	MarketVariationPercentage float64 `json:"regularMarketChangePercent,omitempty"`
	MarketOpen                float64 `json:"regularMarketOpen,omitempty"`
	Symbol                    string  `json:"symbol,omitempty"`
	LongName                  string  `json:"longName,omitempty"`
	ShortName                 string  `json:"shortName,omitempty"`
	Currency                  string  `json:"currency,omitempty"`
}

type PriceMarketTmp struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}
