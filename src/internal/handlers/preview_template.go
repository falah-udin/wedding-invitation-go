package handlers

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/invitations"
)

// ============================================
// PreviewTemplate — GET /preview/template/:slug
// Tampilkan preview template tanpa login, pakai data dummy
// ============================================
func PreviewTemplate(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.String(http.StatusNotFound, "Template tidak ditemukan")
		return
	}

	// Cari template berdasarkan slug
	var template models.Template
	db := database.GetDB()
	if err := db.Where("slug = ? AND is_active = ?", slug, true).
		First(&template).Error; err != nil {
		c.String(http.StatusNotFound, "Template tidak ditemukan")
		return
	}

	log.Printf("═══════════════════════════════════════════════")
	log.Printf("[PREVIEW] START — slug=%s folder=%s", slug, template.Folder)
	log.Printf("═══════════════════════════════════════════════")

	// Bangun data dummy untuk preview
	dummyData := buildDummyData()

	// ============================================
	// DEBUG 1: Cek isi dummyData setelah dibangun
	// ============================================
	debugDummyData(dummyData)

	// Bikin project dummy (tidak disimpan ke DB)
	dummyProject := models.Project{
		ID:           999999,
		Title:        "Preview Undangan",
		Slug:         "preview-" + slug,
		Status:       "published",
		TemplateID:   template.ID,
		Template:     &template,
		DataUndangan: services.ToJSONString(dummyData),
	}

	// ============================================
	// DEBUG 2: Cek DataUndangan JSON (yang bakal dibaca RenderTemplate)
	// ============================================
	log.Printf("[DEBUG-2] DataUndangan JSON length=%d", len(dummyProject.DataUndangan))
	if len(dummyProject.DataUndangan) > 500 {
		log.Printf("[DEBUG-2] DataUndangan preview (500 char): %s...", dummyProject.DataUndangan[:500])
	} else {
		log.Printf("[DEBUG-2] DataUndangan full: %s", dummyProject.DataUndangan)
	}
	log.Printf("[DEBUG-2] Cek 'gallery' ada di JSON? %v", bytes.Contains([]byte(dummyProject.DataUndangan), []byte(`"gallery"`)))
	log.Printf("[DEBUG-2] Cek 'love_stories' ada di JSON? %v", bytes.Contains([]byte(dummyProject.DataUndangan), []byte(`"love_stories"`)))
	log.Printf("[DEBUG-2] Cek 'bank_accounts' ada di JSON? %v", bytes.Contains([]byte(dummyProject.DataUndangan), []byte(`"bank_accounts"`)))

	// Render template
	var buf bytes.Buffer
	err := invitations.RenderTemplate(
		&buf,
		c.Request,
		dummyProject,
		dummyData,
		"Tamu Undangan",
	)
	if err != nil {
		services.LogError("PreviewTemplate.Render", err, map[string]interface{}{
			"slug":   slug,
			"folder": template.Folder,
		})
		c.String(http.StatusInternalServerError, "Gagal render preview: %v", err)
		return
	}

	// ============================================
	// DEBUG 3: Cek HTML hasil render
	// ============================================
	htmlResult := buf.String()
	log.Printf("[DEBUG-3] HTML length=%d", len(htmlResult))
	log.Printf("[DEBUG-3] Cek 'Momen Terindah' ada di HTML? %v", bytes.Contains([]byte(htmlResult), []byte("Momen Terindah")))
	log.Printf("[DEBUG-3] Cek 'gallery-swiper' ada di HTML? %v", bytes.Contains([]byte(htmlResult), []byte("gallery-swiper")))
	log.Printf("[DEBUG-3] Cek 'picsum.photos' ada di HTML? %v", bytes.Contains([]byte(htmlResult), []byte("picsum.photos")))
	log.Printf("[DEBUG-3] Cek 'Perjalanan Cinta' ada di HTML? %v", bytes.Contains([]byte(htmlResult), []byte("Perjalanan Cinta")))
	log.Printf("[DEBUG-3] Cek 'section-gallery' ada di HTML? %v", bytes.Contains([]byte(htmlResult), []byte("section-gallery")))

	// Cari posisi 'section-gallery' kalau ada
	if idx := bytes.Index([]byte(htmlResult), []byte("section-gallery")); idx >= 0 {
		start := idx - 100
		if start < 0 {
			start = 0
		}
		end := idx + 200
		if end > len(htmlResult) {
			end = len(htmlResult)
		}
		log.Printf("[DEBUG-3] Konteks sekitar 'section-gallery': ...%s...", htmlResult[start:end])
	} else {
		log.Printf("[DEBUG-3] ⚠️ 'section-gallery' TIDAK DITEMUKAN di HTML!")
	}

	// Tambah badge "PREVIEW MODE"
	previewBadge := `<div style="position:fixed;top:0;left:0;right:0;background:linear-gradient(90deg,#f59e0b,#d97706);color:#fff;padding:8px 16px;text-align:center;font-family:system-ui,-apple-system,sans-serif;font-size:13px;font-weight:600;z-index:99999;letter-spacing:1px;box-shadow:0 2px 8px rgba(0,0,0,0.15);">
		✦ PREVIEW MODE — Template: <strong>` + template.Name + `</strong> ✦
	</div>`

	htmlResult = injectAfterBody(htmlResult, previewBadge)

	log.Printf("[PREVIEW] END — slug=%s", slug)
	log.Printf("═══════════════════════════════════════════════")

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlResult))
}

// ============================================
// debugDummyData — log isi dummyData untuk cek struktur
// ============================================
func debugDummyData(data map[string]interface{}) {
	log.Printf("[DEBUG-1] ──── Cek struktur dummyData ────")
	log.Printf("[DEBUG-1] Total keys: %d", len(data))

	// Cek gallery
	if g, ok := data["gallery"]; ok {
		log.Printf("[DEBUG-1] gallery type=%T", g)
		if arr, ok := g.([]interface{}); ok {
			log.Printf("[DEBUG-1] gallery len=%d", len(arr))
			for i, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					log.Printf("[DEBUG-1]   gallery[%d] = map{url=%v, caption=%v}", i, m["url"], m["caption"])
				} else {
					log.Printf("[DEBUG-1]   gallery[%d] type=%T value=%v", i, item, item)
				}
			}
		} else {
			log.Printf("[DEBUG-1] ⚠️ gallery BUKAN []interface{} — type=%T", g)
		}
	} else {
		log.Printf("[DEBUG-1] ❌ gallery TIDAK ADA di dummyData")
	}

	// Cek love_stories
	if ls, ok := data["love_stories"]; ok {
		log.Printf("[DEBUG-1] love_stories type=%T", ls)
		if arr, ok := ls.([]interface{}); ok {
			log.Printf("[DEBUG-1] love_stories len=%d", len(arr))
		}
	} else {
		log.Printf("[DEBUG-1] ❌ love_stories TIDAK ADA")
	}

	// Cek bank_accounts
	if ba, ok := data["bank_accounts"]; ok {
		log.Printf("[DEBUG-1] bank_accounts type=%T", ba)
		if arr, ok := ba.([]interface{}); ok {
			log.Printf("[DEBUG-1] bank_accounts len=%d", len(arr))
		}
	} else {
		log.Printf("[DEBUG-1] ❌ bank_accounts TIDAK ADA")
	}

	// Cek foto
	log.Printf("[DEBUG-1] groom_photo=%v", data["groom_photo"])
	log.Printf("[DEBUG-1] bride_photo=%v", data["bride_photo"])
	log.Printf("[DEBUG-1] hero_image=%v", data["hero_image"])
	log.Printf("[DEBUG-1] ────────────────────────────────")
}

// ============================================
// buildDummyData — data contoh untuk preview template
//
// ⚠️ PENTING:
//   - Semua array HARUS []interface{} (bukan []map)
//   - Foto pakai URL eksternal (picsum.photos)
// ============================================
func buildDummyData() map[string]interface{} {
	// ============================================
	// FOTO — pakai picsum.photos
	// ============================================
	groomPhoto := "https://picsum.photos/seed/groom-wedding/500/500"
	bridePhoto := "https://picsum.photos/seed/bride-wedding/500/500"
	fatherPhoto := "https://picsum.photos/seed/father-wedding/400/400"
	motherPhoto := "https://picsum.photos/seed/mother-wedding/400/400"
	heroImage := "https://picsum.photos/seed/hero-wedding/1600/900"

	// ============================================
	// GALLERY — 6 foto, []interface{} of map{url, caption}
	// ============================================
	galleryTitles := []string{"Momen Bahagia", "Kebersamaan", "Cinta Sejati", "Momen Indah", "Bahagia Selalu", "Selamanya"}
	galleryItems := make([]interface{}, 0, len(galleryTitles))
	for i, title := range galleryTitles {
		galleryItems = append(galleryItems, map[string]interface{}{
			"url":     "https://picsum.photos/seed/wedding-gallery-" + itoa(i+1) + "/800/800",
			"caption": title,
		})
	}

	// ============================================
	// RETURN
	// ============================================
	return map[string]interface{}{
		"groom_name": "Falahudin",
		"bride_name": "Zakiyah",

		"akad_date":     "Sabtu, 12 Desember 2026",
		"akad_date_raw": "2026-12-12",
		"akad_time":     "09:00",
		"akad_venue":    "Masjid Agung Al-Azhar, Jakarta",
		"akad_datetime": "2026-12-12T09:00:00+07:00",

		"resepsi_date":     "Sabtu, 12 Desember 2026",
		"resepsi_date_raw": "2026-12-12",
		"resepsi_time":     "11:00",
		"resepsi_venue":    "Gedung Balai Kartini, Jakarta",
		"resepsi_datetime": "2026-12-12T11:00:00+07:00",

		"resepsi_label":        "walimatul_ursy",
		"resepsi_label_custom": "Walimatul 'Ursy",

		"show_dates":         "both",
		"show_venue":         "both",
		"show_bank_accounts": "yes",

		"father_groom": "Bapak H. Abdullah",
		"mother_groom": "Ibu Hj. Fatimah",
		"father_bride": "Bapak H. Mahmud",
		"mother_bride": "Ibu Hj. Aisyah",

		"groom_family_origin": "Yogyakarta",
		"bride_family_origin": "Surakarta",
		"kembar_mayang":       "yes",

		"groom_instagram": "https://instagram.com/ahmadfauzi",
		"bride_instagram": "https://instagram.com/sitiaminah",
		"event_instagram": "https://instagram.com/wedding",
		"whatsapp":        "6281234567890",
		"youtube":         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",

		"maps_url":         "https://maps.google.com/?q=-6.2088,106.8456",
		"maps_url_akad":    "https://maps.google.com/?q=Masjid+Agung+Al-Azhar+Jakarta",
		"maps_url_resepsi": "https://maps.google.com/?q=Balai+Kartini+Jakarta",

		"hero_image":         heroImage,
		"groom_photo":        groomPhoto,
		"bride_photo":        bridePhoto,
		"father_groom_photo": fatherPhoto,
		"mother_groom_photo": motherPhoto,
		"father_bride_photo": fatherPhoto,
		"mother_bride_photo": motherPhoto,

		// LOVE STORIES — []interface{}
		"love_stories": []interface{}{
			map[string]interface{}{
				"title": "Pertemuan",
				"desc":  "Kami pertama bertemu di sebuah acara kampus pada tahun 2020. Saat itu, ada getaran yang tak terlukiskan.",
			},
			map[string]interface{}{
				"title": "Perkenalan",
				"desc":  "Perlahan, hari-hari menjadi lebih indah dengan kehadiran satu sama lain.",
			},
			map[string]interface{}{
				"title": "Pendekatan",
				"desc":  "Dari teman menjadi lebih dari sekadar teman. Setiap momen bersama terasa begitu berharga.",
			},
			map[string]interface{}{
				"title": "Komitmen",
				"desc":  "Janji suci untuk saling menjaga dan mencintai selamanya, dalam suka dan duka.",
			},
		},

		// GALLERY — []interface{}
		"gallery": galleryItems,

		// BANK ACCOUNTS — []interface{}
		"bank_accounts": []interface{}{
			map[string]interface{}{
				"id":             "preview-bank-1",
				"type":           "bank",
				"bank_code":      "bca",
				"bank_name":      "BCA",
				"account_number": "1234567890",
				"account_name":   "Falahudin",
				"icon_type":      "library",
				"icon_path":      "",
			},
			map[string]interface{}{
				"id":             "preview-bank-2",
				"type":           "bank",
				"bank_code":      "mandiri",
				"bank_name":      "Mandiri",
				"account_number": "9876543210",
				"account_name":   "Zakiyah",
				"icon_type":      "library",
				"icon_path":      "",
			},
			map[string]interface{}{
				"id":             "preview-bank-3",
				"type":           "ewallet",
				"bank_code":      "gopay",
				"bank_name":      "GoPay",
				"account_number": "081234567890",
				"account_name":   "Falahudin",
				"icon_type":      "library",
				"icon_path":      "",
			},
			map[string]interface{}{
				"id":             "preview-bank-4",
				"type":           "ewallet",
				"bank_code":      "dana",
				"bank_name":      "DANA",
				"account_number": "081234567890",
				"account_name":   "Zakiyah",
				"icon_type":      "library",
				"icon_path":      "",
			},
		},
	}
}

// ============================================
// injectAfterBody — sisipkan konten setelah tag <body>
// ============================================
func injectAfterBody(html, content string) string {
	idx := -1
	for i := 0; i < len(html)-5; i++ {
		if html[i] == '<' && i+5 <= len(html) && html[i:i+5] == "<body" {
			for j := i + 5; j < len(html); j++ {
				if html[j] == '>' {
					idx = j + 1
					break
				}
			}
			break
		}
	}
	if idx == -1 {
		return content + html
	}
	return html[:idx] + content + html[idx:]
}

// ============================================
// HELPER — itoa
// ============================================
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}