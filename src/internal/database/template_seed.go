package database

import (
	"encoding/json"
	"log"

	"wedding-invitation-go/internal/models"
)

// ============================================
// FIELD LIBRARY — master data semua field
// ============================================
type FieldDef map[string]interface{}

func getFieldLibrary() map[string]FieldDef {
	return map[string]FieldDef{
		// ============================================
		// FOTO & GAMBAR
		// ============================================
		"hero_image": {
			"name":        "hero_image",
			"label":       "Foto Hero / Cover",
			"type":        "file",
			"required":    false,
			"accept":      "image/*",
			"group":       "📸 Foto & Gambar",
			"description": "Foto utama untuk cover undangan (landscape)",
		},
		"groom_photo": {
			"name":        "groom_photo",
			"label":       "Foto Mempelai Pria",
			"type":        "file",
			"required":    true,
			"accept":      "image/*",
			"group":       "📸 Foto & Gambar",
			"description": "Foto mempelai pria (1:1)",
		},
		"bride_photo": {
			"name":        "bride_photo",
			"label":       "Foto Mempelai Wanita",
			"type":        "file",
			"required":    true,
			"accept":      "image/*",
			"group":       "📸 Foto & Gambar",
			"description": "Foto mempelai wanita (1:1)",
		},
		"gallery": {
			"name":        "gallery",
			"label":       "Galeri Foto",
			"type":        "file_multiple",
			"required":    false,
			"accept":      "image/*",
			"group":       "📸 Foto & Gambar",
			"description": "Upload foto pernikahan (max 10, 2MB per foto)",
		},

		// ============================================
		// FOTO KELUARGA
		// ============================================
		"father_groom_photo": {
			"name":        "father_groom_photo",
			"label":       "Foto Ayah Pria",
			"type":        "file",
			"required":    false,
			"accept":      "image/*",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"mother_groom_photo": {
			"name":        "mother_groom_photo",
			"label":       "Foto Ibu Pria",
			"type":        "file",
			"required":    false,
			"accept":      "image/*",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"father_bride_photo": {
			"name":        "father_bride_photo",
			"label":       "Foto Ayah Wanita",
			"type":        "file",
			"required":    false,
			"accept":      "image/*",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"mother_bride_photo": {
			"name":        "mother_bride_photo",
			"label":       "Foto Ibu Wanita",
			"type":        "file",
			"required":    false,
			"accept":      "image/*",
			"group":       "👨‍👩‍👧 Keluarga",
		},

		// ============================================
		// NAMA ORANG TUA
		// ============================================
		"father_groom": {
			"name":        "father_groom",
			"label":       "Nama Ayah Mempelai Pria",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Bapak H. Abdullah",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"mother_groom": {
			"name":        "mother_groom",
			"label":       "Nama Ibu Mempelai Pria",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Ibu Hj. Fatimah",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"father_bride": {
			"name":        "father_bride",
			"label":       "Nama Ayah Mempelai Wanita",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Bapak H. Mahmud",
			"group":       "👨‍👩‍👧 Keluarga",
		},
		"mother_bride": {
			"name":        "mother_bride",
			"label":       "Nama Ibu Mempelai Wanita",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Ibu Hj. Aisyah",
			"group":       "👨‍👩‍👧 Keluarga",
		},

		// ============================================
		// ASAL & TRADISI
		// ============================================
		"groom_family_origin": {
			"name":        "groom_family_origin",
			"label":       "Asal Keluarga Pria",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Yogyakarta",
			"group":       "🏛️ Asal & Tradisi",
		},
		"bride_family_origin": {
			"name":        "bride_family_origin",
			"label":       "Asal Keluarga Wanita",
			"type":        "text",
			"required":    false,
			"placeholder": "Contoh: Surakarta",
			"group":       "🏛️ Asal & Tradisi",
		},
		"kembar_mayang": {
			"name":        "kembar_mayang",
			"label":       "Kembar Mayang",
			"type":        "select",
			"required":    false,
			"options": map[string]string{
				"yes": "Ya, Pakai Kembar Mayang",
				"no":  "Tidak Pakai",
			},
			"group":       "🏛️ Asal & Tradisi",
			"description": "Kembar Mayang adalah simbol adat Jawa",
		},

		// ============================================
		// SOSIAL MEDIA
		// ============================================
		"groom_instagram": {
			"name":        "groom_instagram",
			"label":       "Instagram Mempelai Pria",
			"type":        "url",
			"required":    false,
			"placeholder": "https://instagram.com/username",
			"group":       "🌐 Sosial Media",
		},
		"bride_instagram": {
			"name":        "bride_instagram",
			"label":       "Instagram Mempelai Wanita",
			"type":        "url",
			"required":    false,
			"placeholder": "https://instagram.com/username",
			"group":       "🌐 Sosial Media",
		},
		"event_instagram": {
			"name":        "event_instagram",
			"label":       "Instagram Acara",
			"type":        "url",
			"required":    false,
			"placeholder": "https://instagram.com/wedding",
			"group":       "🌐 Sosial Media",
		},
		"whatsapp": {
			"name":        "whatsapp",
			"label":       "WhatsApp Acara",
			"type":        "text",
			"required":    false,
			"placeholder": "6281234567890",
			"group":       "🌐 Sosial Media",
		},
		"youtube": {
			"name":        "youtube",
			"label":       "YouTube Acara",
			"type":        "url",
			"required":    false,
			"placeholder": "https://youtube.com/...",
			"group":       "🌐 Sosial Media",
		},

		// ============================================
		// LOKASI / MAPS
		// ============================================
		"maps_url_akad": {
			"name":        "maps_url_akad",
			"label":       "Google Maps Lokasi Akad",
			"type":        "url",
			"required":    false,
			"placeholder": "https://maps.google.com/...",
			"group":       "📍 Lokasi",
		},
		"maps_url_resepsi": {
			"name":        "maps_url_resepsi",
			"label":       "Google Maps Lokasi Resepsi",
			"type":        "url",
			"required":    false,
			"placeholder": "https://maps.google.com/...",
			"group":       "📍 Lokasi",
		},

		// ============================================
		// CERITA CINTA
		// ============================================
		"love_stories": {
			"name":     "love_stories",
			"label":    "Cerita Cinta",
			"type":     "repeater",
			"required": false,
			"group":    "❤️ Cerita Cinta",
			"fields": []map[string]interface{}{
				{"name": "title", "label": "Judul", "type": "text", "placeholder": "Contoh: Pertemuan"},
				{"name": "desc", "label": "Deskripsi", "type": "textarea", "placeholder": "Ceritakan momen...", "rows": 2},
			},
			"default_count": 4,
			"default_stories": []map[string]string{
				{"title": "Pertemuan", "desc": "Saat pertama kali bertemu, ada getaran yang tak terlukiskan..."},
				{"title": "Perkenalan", "desc": "Perlahan, hari-hari menjadi lebih indah dengan kehadiranmu..."},
				{"title": "Pendekatan", "desc": "Dari teman menjadi lebih dari sekadar teman..."},
				{"title": "Komitmen", "desc": "Janji suci untuk saling menjaga dan mencintai selamanya..."},
			},
			"description": "Tambahkan cerita perjalanan cinta Anda (min 1, max 10)",
		},
	}
}

// ============================================
// DEFAULT FIELDS — urutan standar
// ============================================
func getDefaultFields() []string {
	return []string{
		"hero_image",
		"groom_photo",
		"bride_photo",
		"gallery",
		"father_groom",
		"mother_groom",
		"father_bride",
		"mother_bride",
		"groom_instagram",
		"bride_instagram",
		"event_instagram",
		"whatsapp",
		"maps_url_akad",
		"maps_url_resepsi",
		"love_stories",
	}
}

// ============================================
// BUILD SCHEMA — dari list field name
// ============================================
func buildSchema(fieldNames []string, library map[string]FieldDef) []map[string]interface{} {
	var schema []map[string]interface{}
	for _, name := range fieldNames {
		if field, ok := library[name]; ok {
			schema = append(schema, map[string]interface{}(field))
		}
	}
	return schema
}

// ============================================
// INSERT / UPDATE TEMPLATE
// ============================================
func upsertTemplate(name, slug, folder, description string, schema []map[string]interface{}, order int) error {
	schemaJSON, _ := json.Marshal(schema)

	var existing models.Template
	result := DB.Where("slug = ?", slug).First(&existing)

	data := models.Template{
		Name:         name,
		Slug:         slug,
		Folder:       folder,
		Sections:     `["cover","couple","love-story","event","gallery","rsvp"]`,
		Description:  &description,
		Fields:       `[]`,
		FieldsSchema: string(schemaJSON),
		IsActive:     true,
		Order:        order,
	}

	if result.Error == nil {
		// Update
		existing.Name = data.Name
		existing.Folder = data.Folder
		existing.Sections = data.Sections
		existing.Description = data.Description
		existing.Fields = data.Fields
		existing.FieldsSchema = data.FieldsSchema
		existing.IsActive = data.IsActive
		existing.Order = data.Order
		return DB.Save(&existing).Error
	}

	// Create
	return DB.Create(&data).Error
}

// ============================================
// SEED TEMPLATES — main function
// ============================================
func seedTemplates() error {
	library := getFieldLibrary()
	defaultFields := getDefaultFields()

	// ============================================
	// 1. RUSTIC WOOD (default)
	// ============================================
	if err := upsertTemplate(
		"Rustic Wood",
		"rustic-wood",
		"rustic-wood",
		"Tema kayu alami dengan nuansa hangat dan natural. Cocok untuk pernikahan outdoor.",
		buildSchema(defaultFields, library),
		1,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Rustic Wood")

	// ============================================
	// 2. MODERN MINIMALIST (default + youtube)
	// ============================================
	modernFields := []string{}
	for _, f := range defaultFields {
		modernFields = append(modernFields, f)
		if f == "whatsapp" {
			modernFields = append(modernFields, "youtube")
		}
	}
	if err := upsertTemplate(
		"Modern Minimalist",
		"modern-minimalist",
		"modern-minimalist",
		"Desain simpel elegan dengan gaya modern minimalis. Cocok untuk konsep bersih dan sophisticated.",
		buildSchema(modernFields, library),
		2,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Modern Minimalist")

	// ============================================
	// 3. ELEGANT GOLD (default)
	// ============================================
	if err := upsertTemplate(
		"Elegant Gold",
		"elegant-gold",
		"elegant-gold",
		"Template elegan dengan nuansa emas mewah. Dilengkapi efek shimmer dan ornamen klasik.",
		buildSchema(defaultFields, library),
		3,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Elegant Gold")

	// ============================================
	// 4. TRADITIONAL JAVA (default + adat jawa)
	// ============================================
	javaFields := []string{}
	for _, f := range defaultFields {
		javaFields = append(javaFields, f)
		if f == "mother_bride" {
			javaFields = append(javaFields, "groom_family_origin", "bride_family_origin", "kembar_mayang")
		}
	}
	if err := upsertTemplate(
		"Traditional Java",
		"traditional-java",
		"traditional-java",
		"Tema adat jawa kental dengan budaya dan tradisi. Dilengkapi motif batik dan ornamen khas Jawa.",
		buildSchema(javaFields, library),
		4,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Traditional Java")

	// ============================================
	// 5. MUSLIM ELEGAN (default tanpa hero + foto ortu)
	// ============================================
	muslimFields := []string{}
	for _, f := range defaultFields {
		if f == "hero_image" {
			continue
		}
		muslimFields = append(muslimFields, f)
		if f == "gallery" {
			muslimFields = append(muslimFields,
				"father_groom_photo",
				"mother_groom_photo",
				"father_bride_photo",
				"mother_bride_photo",
			)
		}
	}
	if err := upsertTemplate(
		"Muslim Elegan",
		"muslim-elegan",       
		"muslim-elegan",      
		"Template dengan nuansa islami yang elegan, dilengkapi ayat Al-Quran dan desain hangat.",
		buildSchema(muslimFields, library),
		5,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Muslim Elegan")

	return nil
}
