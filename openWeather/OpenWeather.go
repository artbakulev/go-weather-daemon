package openWeather

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const BaseUrl = "https://api.openweathermap.org/data/2.5/weather"
const AbsoluteZero = 273.15

type API struct {
	AppId   string
	Timeout time.Duration
}

//func NewOpenWeatherAPI(appId string) *API {
//	api := API{AppId: appId}
//	return &api
//}

func (v API) buildUrl(query string) (string, error) {
	base, err := url.Parse(BaseUrl)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("q", query)
	params.Add("appid", v.AppId)
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (v API) GetWeatherForQuery(ctx context.Context, query string) (float32, error) {
	ctx, _ = context.WithTimeout(ctx, v.Timeout)
	u, err := v.buildUrl(query)
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, err
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	weatherData := WeatherData{}
	err = json.Unmarshal(buffer, &weatherData)
	if err != nil {
		return 0.0, err
	}
	return weatherData.Main.Temperature - AbsoluteZero, nil
}
