package models

// BankInfo — info bank/e-wallet
type BankInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"` // bank, ewallet, qris, other
	Full string `json:"full"`
}

// BankAccount — struktur 1 rekening yang disimpan di data_undangan
type BankAccount struct {
	ID            string  `json:"id"`
	Type          string  `json:"type"`           // bank, ewallet, qris, other
	BankCode      string  `json:"bank_code"`      // bca, mandiri, gopay, dll
	BankName      string  `json:"bank_name"`      // custom name (kalau other)
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	IconType      string  `json:"icon_type"`      // library, upload
	IconPath      *string `json:"icon_path,omitempty"`
}

// GetBankList — daftar bank & e-wallet yang tersedia
// WAJIB SAMA dengan BANK_LIST di JavaScript (general_script.templ)
func GetBankList() []BankInfo {
	return []BankInfo{
		// Bank
		{Code: "bca", Name: "BCA", Type: "bank", Full: "Bank Central Asia"},
		{Code: "mandiri", Name: "Mandiri", Type: "bank", Full: "Bank Mandiri"},
		{Code: "bni", Name: "BNI", Type: "bank", Full: "Bank Negara Indonesia"},
		{Code: "bri", Name: "BRI", Type: "bank", Full: "Bank Rakyat Indonesia"},
		{Code: "bsi", Name: "BSI", Type: "bank", Full: "Bank Syariah Indonesia"},
		{Code: "cimb", Name: "CIMB Niaga", Type: "bank", Full: "CIMB Niaga"},
		{Code: "permata", Name: "Permata", Type: "bank", Full: "Bank Permata"},
		{Code: "danamon", Name: "Danamon", Type: "bank", Full: "Bank Danamon"},
		{Code: "btn", Name: "BTN", Type: "bank", Full: "Bank Tabungan Negara"},
		{Code: "mega", Name: "Mega", Type: "bank", Full: "Bank Mega"},
		{Code: "bjb", Name: "BJB", Type: "bank", Full: "Bank Jabar Banten"},
		{Code: "jago", Name: "Jago", Type: "bank", Full: "Bank Jago"},
		{Code: "seabank", Name: "SeaBank", Type: "bank", Full: "SeaBank Indonesia"},
		{Code: "blu", Name: "Blu BCA", Type: "bank", Full: "Blu by BCA Digital"},
		{Code: "jenius", Name: "Jenius", Type: "bank", Full: "Jenius BTPN"},

		// E-Wallet
		{Code: "gopay", Name: "GoPay", Type: "ewallet", Full: "GoPay"},
		{Code: "ovo", Name: "OVO", Type: "ewallet", Full: "OVO"},
		{Code: "dana", Name: "DANA", Type: "ewallet", Full: "DANA"},
		{Code: "shopeepay", Name: "ShopeePay", Type: "ewallet", Full: "ShopeePay"},
		{Code: "linkaja", Name: "LinkAja", Type: "ewallet", Full: "LinkAja"},

		// QRIS & Other
		{Code: "qris", Name: "QRIS", Type: "qris", Full: "QRIS (Semua Bank/E-Wallet)"},
		{Code: "other", Name: "Lainnya", Type: "other", Full: "Lainnya"},
	}
}

// GetBankName — ambil nama bank berdasarkan code
func GetBankName(code string) string {
	for _, b := range GetBankList() {
		if b.Code == code {
			return b.Name
		}
	}
	return code
}
