package invitation

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/invitations"
)

// ============================================
// ShowInvitation — GET /invitation/:slug
// Tampilkan undangan ke tamu (publik)
// ============================================
func ShowInvitation(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.String(http.StatusNotFound, "Undangan tidak ditemukan")
		return
	}

	// Cari project berdasarkan slug
	var project models.Project
	db := database.GetDB()
	if err := db.Preload("Template").Preload("User").
		Where("slug = ?", slug).
		First(&project).Error; err != nil {
		c.String(http.StatusNotFound, "Undangan tidak ditemukan")
		return
	}

	// Cek status: hanya published yang bisa diakses
	if project.Status != "published" {
		c.String(http.StatusNotFound, "Undangan belum dipublikasikan")
		return
	}

	// Increment views
	db.Model(&project).UpdateColumn("total_views", project.TotalViews+1)

	// Ambil data undangan (merge general + specific)
	generalData := services.ParseJSONMap(project.DataUndangan)
	specificData := services.ParseJSONMap(project.TemplateSpecificData)

	allData := make(map[string]interface{})
	for k, v := range generalData {
		allData[k] = v
	}
	for k, v := range specificData {
		allData[k] = v
	}

	// === INJECT DATA TAMBAHAN UNTUK TEMPLATE ===
	// 1. Site config (untuk footer web)
	siteCfg := services.GetSiteSetting()
	allData["_site_config"] = siteCfg

	// 2. Bank list (untuk lookup logo bank di amplop digital)
	bankList := make([]map[string]string, 0)
	for _, b := range models.GetBankList() {
		bankList = append(bankList, map[string]string{
			"code": b.Code,
			"name": b.Name,
			"type": b.Type,
			"full": b.Full,
		})
	}
	allData["_bank_list"] = bankList

	// Nama tamu dari query parameter ?to=Nama
	guestName := c.Query("to")
	if guestName == "" {
		guestName = "Tamu Undangan"
	}
	// Normalize: ganti underscore/+ jadi spasi
	guestName = strings.ReplaceAll(guestName, "_", " ")
	guestName = strings.ReplaceAll(guestName, "+", " ")
	guestName = strings.Join(strings.Fields(guestName), " ")

	// 3. Existing RSVP (untuk pre-fill form)
	if guestName != "Tamu Undangan" {
		var existingRsvp models.Rsvp
		err := db.Where("project_id = ? AND LOWER(guest_name) = ?",
			project.ID, strings.ToLower(guestName)).
			First(&existingRsvp).Error
		if err == nil {
			allData["_existing_rsvp"] = map[string]interface{}{
				"attendance":   existingRsvp.Attendance,
				"total_guests": existingRsvp.TotalGuests,
				"message":      existingRsvp.GetMessage(),
			}
		}
	}

	// Render template sesuai folder
	// Template view: views/invitations/{folder}/index.templ
	var buf bytes.Buffer
	err := invitations.RenderTemplate(
		&buf,
		c.Request,
		project,
		allData,
		guestName,
	)

	if err != nil {
		services.LogError("ShowInvitation.Render", err, map[string]interface{}{
			"slug":   slug,
			"folder": project.Template.Folder,
		})
		c.String(http.StatusInternalServerError, "Gagal render undangan: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// RsvpSubmit — POST /invitation/:slug/rsvp
// Submit RSVP dari tamu
// ============================================
func RsvpSubmit(c *gin.Context) {
	slug := c.Param("slug")

	var project models.Project
	db := database.GetDB()
	if err := db.Where("slug = ?", slug).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Undangan tidak ditemukan"})
		return
	}

	guestName := strings.TrimSpace(c.PostForm("guest_name"))
	attendance := strings.TrimSpace(c.PostForm("attendance"))
	totalGuests := c.PostForm("total_guests")
	message := strings.TrimSpace(c.PostForm("message"))

	if guestName == "" || attendance == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Nama dan kehadiran wajib diisi"})
		return
	}

	// Parse total guests
	total := 1
	if totalGuests != "" {
		for _, ch := range totalGuests {
			if ch >= '0' && ch <= '9' {
				total = total*10 + int(ch-'0')
			}
		}
	}

	// Normalize nama
	guestName = strings.Join(strings.Fields(guestName), " ")

	// Cek duplikat
	var existing models.Rsvp
	err := db.Where("project_id = ? AND LOWER(guest_name) = ?", project.ID, strings.ToLower(guestName)).
		First(&existing).Error

	if err == nil {
		// Update existing
		existing.Attendance = attendance
		existing.TotalGuests = total
		if message != "" {
			existing.Message = &message
		}
		db.Save(&existing)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Konfirmasi berhasil diperbarui!",
			"updated": true,
			"data":    existing,
		})
		return
	}

	// Buat RSVP baru
	rsvp := models.Rsvp{
		ProjectID:   project.ID,
		GuestName:   guestName,
		Attendance:  attendance,
		TotalGuests: total,
	}
	if message != "" {
		rsvp.Message = &message
	}

	if err := db.Create(&rsvp).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menyimpan RSVP"})
		return
	}

	// Increment total_rsvp di project
	db.Model(&project).UpdateColumn("total_rsvp", project.TotalRsvp+1)

	// Sync ke InvitationGuest
	var existingGuest models.InvitationGuest
	if err := db.Where("project_id = ? AND LOWER(name) = ?", project.ID, strings.ToLower(guestName)).
		First(&existingGuest).Error; err == nil {
		existingGuest.RsvpID = &rsvp.ID
		db.Save(&existingGuest)
	} else {
		guest := models.InvitationGuest{
			ProjectID: project.ID,
			RsvpID:    &rsvp.ID,
			Name:      guestName,
		}
		db.Create(&guest)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Terima kasih! Konfirmasi kehadiran berhasil dikirim.",
		"updated": false,
		"data":    rsvp,
	})
}

// ============================================
// RsvpList — GET /invitation/:slug/rsvp-list
// Ambil daftar RSVP (JSON untuk AJAX)
// ============================================
func RsvpList(c *gin.Context) {
	slug := c.Param("slug")

	var project models.Project
	db := database.GetDB()
	if err := db.Where("slug = ?", slug).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	var rsvps []models.Rsvp
	db.Where("project_id = ?", project.ID).
		Order("created_at DESC").
		Limit(50).
		Find(&rsvps)

	c.JSON(http.StatusOK, rsvps)
}
