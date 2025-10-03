package dto

import "github.com/biraneves/fc-labs-weather/internal/domain/entity"

type RequestInDto struct {
	CEP entity.Cep `json:"cep"`
}

type RequestOutDto struct {
	TempC entity.TemperatureCelsius    `json:"temp_C"`
	TempF entity.TemperatureFahrenheit `json:"temp_F"`
	TempK entity.TemperatureKelvin     `json:"temp_K"`
}

type ViaCEPRequestDto struct {
	CEP entity.Cep `json:"cep"`
}

type ViaCEPResponseDto struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type WeatherAPIRequestDto struct {
	Q string `json:"q"`
}

type WeatherAPIResponseDto struct {
	Location struct {
		Name          string  `json:"name"`
		Region        string  `json:"region"`
		Country       string  `json:"country"`
		Lat           float64 `json:"lat"`
		Lon           float64 `json:"lon"`
		TzId          string  `json:"tz_id"`
		LocationEpoch int     `json:"location_epoch"`
		Localtime     string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
		FeelsLikeC float64 `json:"feelslike_c"`
		FeelsLikeF float64 `json:"feelskile_f"`
		WindChillC float64 `json:"windchill_c"`
		WindChillF float64 `json:"windchill_f"`
		HeatIndexC float64 `json:"heatindex_c"`
		HeatIndexF float64 `json:"heatindex_f"`
		DewPointC  float64 `json:"dewpoint_c"`
		DewPointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}
