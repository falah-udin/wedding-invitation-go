package models

// Template — model sederhana untuk template undangan
// Nanti field-nya akan di-map dari database MySQL (tabel: templates)
type Template struct {
	ID          uint
	Name        string
	Slug        string
	Folder      string
	Thumbnail   string
	Description string
	IsActive    bool
	Order       int
}

// DummyTemplates — data dummy untuk test UI
// Setelah database siap, ini diganti query GORM ke tabel templates
func DummyTemplates() []Template {
	return []Template{
		{
			ID:          1,
			Name:        "Rustic Wood",
			Slug:        "rustic-wood",
			Folder:      "rustic-wood",
			Thumbnail:   "/images/templates/rustic-wood.jpg",
			Description: "Tema kayu alami dengan nuansa hangat dan natural.",
			IsActive:    true,
			Order:       1,
		},
		{
			ID:          2,
			Name:        "Modern Minimalist",
			Slug:        "modern-minimalist",
			Folder:      "modern-minimalist",
			Thumbnail:   "/images/templates/modern-minimalist.jpg",
			Description: "Desain simpel elegan dengan gaya modern minimalis.",
			IsActive:    true,
			Order:       2,
		},
		{
			ID:          3,
			Name:        "Elegant Gold",
			Slug:        "elegant-gold",
			Folder:      "elegant-gold",
			Thumbnail:   "/images/templates/elegant-gold.jpg",
			Description: "Template elegan dengan nuansa emas mewah.",
			IsActive:    true,
			Order:       3,
		},
		{
			ID:          4,
			Name:        "Traditional Java",
			Slug:        "traditional-java",
			Folder:      "traditional-java",
			Thumbnail:   "/images/templates/traditional-java.jpg",
			Description: "Tema adat jawa yang kental dengan budaya dan tradisi.",
			IsActive:    true,
			Order:       4,
		},
		{
			ID:          5,
			Name:        "Muslim Elegan",
			Slug:        "santri-islami",
			Folder:      "santri-islami",
			Thumbnail:   "/images/templates/santri-islami.jpg",
			Description: "Template islami dengan nuansa elegan dan ayat Al-Quran.",
			IsActive:    true,
			Order:       5,
		},
	}
}