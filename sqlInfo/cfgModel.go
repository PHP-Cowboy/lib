package sqlInfo

type GlobalPersonalCfgItem struct {
	Id              int     `json:"Id"`              //
	MinRecharge     int     `json:"MinRecharge"`     //
	WithDrawRate    float64 `json:"WithDrawRate"`    //
	BaseWinFactor   float64 `json:"BaseWinFactor"`   //
	WinLineRate     float64 `json:"WinLineRate"`     //
	AddFactor       float64 `json:"AddFactor"`       //
	LoseBeginRate   float64 `json:"LoseBeginRate"`   //
	LoseLineRate    float64 `json:"LoseLineRate"`    //
	UpLimitRate     float64 `json:"UpLimitRate"`     //
	MaxBaseRate     float64 `json:"MaxBaseRate"`     //
	MaxWinRate      float64 `json:"MaxWinRate"`      //
	DailyWin        int     `json:"DailyWin"`        //
	DailyWinValue   int64   `json:"DailyWinValue"`   //
	SaveTime        int     `json:"SaveTime"`        //
	BigWinRate      int64   `json:"BigWinRate"`      //
	BigRateInterval []int   `json:"BigRateInterval"` //
	MinWBetRate     float64 `json:"MinWBetRate"`     //
	MaxWBetRate     float64 `json:"MaxWBetRate"`     //
	WBetFactor      float64 `json:"WBetFactor"`      //
	WBetRate        int     `json:"WBetRate"`        //

}

func (p *GlobalPersonalCfgItem) GetUpLimitRate() float64 {
	return p.UpLimitRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetLoseBeginRate() float64 {
	return p.LoseBeginRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetAddFactor() float64 {
	return p.AddFactor / 1000.0
}
func (p *GlobalPersonalCfgItem) GetWithDrawRate() float64 {
	return p.WithDrawRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetLoseLineRate() float64 {
	return p.LoseLineRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetBaseWinFactor() float64 {
	return p.BaseWinFactor / 1000.0
}
func (p *GlobalPersonalCfgItem) GetWinLineRate() float64 {
	return p.WinLineRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetMaxBaseRate() float64 {
	return p.MaxBaseRate / 1000.0
}
func (p *GlobalPersonalCfgItem) GetMaxWinRate() float64 {
	return p.MaxWinRate / 1000.0
}

type GlobalPersonalCfg struct {
	Items []GlobalPersonalCfgItem
}

// /////////
func (p *GlobalPersonalCfg) GetGlobalPersonalCfgItem(rec int64) *GlobalPersonalCfgItem {
	for i := len(p.Items) - 1; i >= 0; i-- {
		if rec >= int64(p.Items[i].MinRecharge) {
			return &p.Items[i]
		}
	}
	return nil
}

func (p *GlobalPersonalCfg) GetDailyWin(rec int64) float64 {
	if v := p.GetGlobalPersonalCfgItem(rec); v != nil {
		return float64(v.DailyWin) / 1000.0
	}
	return 0
}
func (p *GlobalPersonalCfg) GetUpLimitRate(rec int64) float64 {
	if v := p.GetGlobalPersonalCfgItem(rec); v != nil {
		return float64(v.UpLimitRate) / 1000.0
	}
	return 0
}
