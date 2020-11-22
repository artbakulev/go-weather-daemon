package openWeather

type WeatherData struct {
	Main struct {
		Temperature float32 `json:"temp"`
	} `json:"main"`
}
