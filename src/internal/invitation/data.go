package invitation

import (
	"wedding-invitation-go/internal/models"
)

// TemplateData — semua data yang dibutuhkan template undangan
type TemplateData struct {
	// Mempelai
	GroomName string
	BrideName string

	// Foto
	GroomPhoto string
	BridePhoto string
	HeroImage  string

	// Foto keluarga
	FatherGroomPhoto string
	MotherGroomPhoto string
	FatherBridePhoto string
	MotherBridePhoto string

	// Orang tua
	FatherGroom string
	MotherGroom string
	FatherBride string
	MotherBride string

	// Asal keluarga
	GroomFamilyOrigin string
	BrideFamilyOrigin string
	KembarMayang      string

	// Sosial media
	GroomInstagram string
	BrideInstagram string
	EventInstagram string
	Whatsapp       string
	Youtube        string

	// Akad
	AkadDate  string
	AkadTime  string
	AkadVenue string

	// Resepsi
	ResepsiDate        string
	ResepsiTime        string
	ResepsiVenue       string
	ResepsiLabel       string
	ResepsiLabelCustom string

	// Maps
	MapsURLAkad    string
	MapsURLResepsi string

	// Tampilan cover
	ShowDates        string
	ShowVenue        string
	ShowBankAccounts string

	// Parsed
	BankAccounts []BankAccount
	Gallery      []string
	LoveStories  []LoveStory

	// Runtime inject
	SiteConfig   models.Config
	BankList     []models.BankInfo
	ExistingRsvp *ExistingRsvp

	// Dinamis
	Specific map[string]interface{}

	// Metadata
	Project      models.Project
	GuestName    string
	IsNamedGuest bool
	MusicURL     string
}

type BankAccount struct {
	ID            string
	Type          string
	BankCode      string
	BankName      string
	AccountNumber string
	AccountName   string
	IconType      string
	IconPath      string
}

type LoveStory struct {
	Title string
	Desc  string
}

type ExistingRsvp struct {
	Attendance  string
	TotalGuests int
	Message     string
}
