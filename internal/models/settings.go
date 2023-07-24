package models

// AppSettings is a struct that holds the settings for the application
type AppSettings struct {
	Enabled         bool    `json:"enabled"`
	Mode            string  `json:"mode"`
	ShowCPUTemp     bool    `json:"show_cpu_temp"`
	CPUFanSpeed     string  `json:"cpu_fan_speed"`
	CPULoadOrRAMUse string  `json:"cpu_load_or_ram_use"`
	ShowSongTitle   bool    `json:"show_song_title"`
	ShowEqualizer   bool    `json:"show_equalizer"`
	EqualizerDemo   bool    `json:"equalizer_demo"`
	Brightness      float64 `json:"brightness"`
	IsMirrored      bool    `json:"is_mirrored"`
}
