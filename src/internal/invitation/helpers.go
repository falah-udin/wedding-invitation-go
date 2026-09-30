package invitation

import (
	"strconv"
	"time"
)

const (
	DefaultGroomPhotoPath  = "/storage/defaults/groom.svg"
	DefaultBridePhotoPath  = "/storage/defaults/bride.svg"
	DefaultFatherPhotoPath = "/storage/defaults/father.svg"
	DefaultMotherPhotoPath = "/storage/defaults/mother.svg"
	DefaultPhotoPath       = "/storage/defaults/photo.svg"
)

func (d *TemplateData) GroomPhotoOrDefault() string {
	if d.GroomPhoto != "" {
		return d.GroomPhoto
	}
	return DefaultGroomPhotoPath
}

func (d *TemplateData) BridePhotoOrDefault() string {
	if d.BridePhoto != "" {
		return d.BridePhoto
	}
	return DefaultBridePhotoPath
}

func (d *TemplateData) FatherGroomPhotoOrDefault() string {
	if d.FatherGroomPhoto != "" {
		return d.FatherGroomPhoto
	}
	return DefaultFatherPhotoPath
}

func (d *TemplateData) MotherGroomPhotoOrDefault() string {
	if d.MotherGroomPhoto != "" {
		return d.MotherGroomPhoto
	}
	return DefaultMotherPhotoPath
}

func (d *TemplateData) FatherBridePhotoOrDefault() string {
	if d.FatherBridePhoto != "" {
		return d.FatherBridePhoto
	}
	return DefaultFatherPhotoPath
}

func (d *TemplateData) MotherBridePhotoOrDefault() string {
	if d.MotherBridePhoto != "" {
		return d.MotherBridePhoto
	}
	return DefaultMotherPhotoPath
}

func (d *TemplateData) HeroImageOr(fallback string) string {
	if d.HeroImage != "" {
		return d.HeroImage
	}
	return fallback
}

func (d *TemplateData) HasHeroImage() bool {
	return d.HeroImage != ""
}

func (d *TemplateData) ResepsiLabelDisplay() string {
	switch d.ResepsiLabel {
	case "walimatul_ursy":
		return "Walimatul Ursy"
	case "walimah":
		return "Walimah"
	case "custom":
		if d.ResepsiLabelCustom != "" {
			return d.ResepsiLabelCustom
		}
		return "Resepsi"
	default:
		return "Resepsi"
	}
}

func (d *TemplateData) ResepsiArabic() string {
	switch d.ResepsiLabel {
	case "walimatul_ursy":
		return "وَلِيمَةُ الْعُرْسِ"
	case "walimah":
		return "وَلِيمَة"
	}
	return ""
}

func (d *TemplateData) IsIslamicResepsi() bool {
	return d.ResepsiLabel == "walimatul_ursy" || d.ResepsiLabel == "walimah"
}

func (d *TemplateData) ParentsGroom() string {
	return joinParents(d.FatherGroom, d.MotherGroom)
}

func (d *TemplateData) ParentsBride() string {
	return joinParents(d.FatherBride, d.MotherBride)
}

func joinParents(father, mother string) string {
	if father != "" && mother != "" {
		return father + " & " + mother
	}
	if father != "" {
		return father
	}
	if mother != "" {
		return mother
	}
	return "Bapak & Ibu"
}

func (d *TemplateData) ShouldShowBankAccounts() bool {
	return d.ShowBankAccounts == "yes"
}

func (d *TemplateData) ShouldShowAkadDate() bool {
	return d.ShowDates == "both" || d.ShowDates == "akad"
}

func (d *TemplateData) ShouldShowResepsiDate() bool {
	return d.ShowDates == "both" || d.ShowDates == "resepsi"
}

func (d *TemplateData) ShouldShowBothDates() bool {
	return d.ShowDates == "both"
}

func (d *TemplateData) ShouldShowAkadVenue() bool {
	return d.ShowVenue == "both" || d.ShowVenue == "akad"
}

func (d *TemplateData) ShouldShowResepsiVenue() bool {
	return d.ShowVenue == "both" || d.ShowVenue == "resepsi"
}

func (d *TemplateData) ShouldShowBothVenues() bool {
	return d.ShowVenue == "both"
}

func (d *TemplateData) HasExistingRsvp() bool {
	return d.ExistingRsvp != nil
}

func (d *TemplateData) ExistingRsvpField(key, def string) string {
	if d.ExistingRsvp == nil {
		return def
	}
	switch key {
	case "attendance":
		if d.ExistingRsvp.Attendance != "" {
			return d.ExistingRsvp.Attendance
		}
	case "message":
		if d.ExistingRsvp.Message != "" {
			return d.ExistingRsvp.Message
		}
	}
	return def
}

func (d *TemplateData) ExistingRsvpGuests() int {
	if d.ExistingRsvp == nil || d.ExistingRsvp.TotalGuests <= 0 {
		return 1
	}
	return d.ExistingRsvp.TotalGuests
}

func GetBankTypeBadge(acc BankAccount) string {
	switch acc.Type {
	case "bank":
		return "BANK"
	case "ewallet":
		return "E-WALLET"
	case "qris":
		return "QRIS"
	default:
		return "LAINNYA"
	}
}

func GetAccountBankName(acc BankAccount) string {
	if acc.BankName != "" {
		return acc.BankName
	}
	if name := getBankNameFromLibrary(acc.BankCode); name != "" {
		return name
	}
	return acc.BankCode
}

func GetAccountInitial(acc BankAccount) string {
	name := GetAccountBankName(acc)
	if len(name) >= 2 {
		return name[:2]
	}
	return name
}

func GetAccountNumber(acc BankAccount) string {
	if acc.AccountNumber == "" {
		return "-"
	}
	return acc.AccountNumber
}

func GetAccountName(acc BankAccount) string {
	if acc.AccountName == "" {
		return "-"
	}
	return acc.AccountName
}

func HasCustomIcon(acc BankAccount) bool {
	return acc.IconType == "upload" && acc.IconPath != ""
}

func GetAccountIconURL(acc BankAccount) string {
	if acc.IconPath == "" {
		return ""
	}
	if len(acc.IconPath) > 0 && acc.IconPath[0] != '/' && acc.IconPath[0] != 'h' {
		return "/storage/" + acc.IconPath
	}
	return acc.IconPath
}

func getBankNameFromLibrary(code string) string {
	bankMap := map[string]string{
		"bca": "BCA", "mandiri": "Mandiri", "bni": "BNI", "bri": "BRI",
		"bsi": "BSI", "cimb": "CIMB Niaga", "permata": "Permata",
		"danamon": "Danamon", "btn": "BTN", "mega": "Mega", "bjb": "BJB",
		"jago": "Jago", "seabank": "SeaBank", "blu": "Blu BCA", "jenius": "Jenius",
		"gopay": "GoPay", "ovo": "OVO", "dana": "DANA", "shopeepay": "ShopeePay",
		"linkaja": "LinkAja", "qris": "QRIS", "other": "Lainnya",
	}
	return bankMap[code]
}

func FormatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

func FormatInt(n int) string {
	return strconv.Itoa(n)
}

func CurrentYear() int {
	return time.Now().Year()
}

func GuestGreeting(isNamed bool) string {
	if isNamed {
		return "Kepada"
	}
	return "Salam Hangat untuk"
}
