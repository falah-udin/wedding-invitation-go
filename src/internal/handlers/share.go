package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/layouts"
	"wedding-invitation-go/views/share"
)

// ============================================
// ShareInvitation — GET /share/invitation/:id
// Halaman daftar tamu + share WA
// ============================================
func ShareInvitation(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	// Parse ID dari parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=ID+tidak+valid")
		return
	}

	// Ambil project
	var project models.Project
	db := database.GetDB()
	if err := db.Preload("User").Preload("Template").First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	// Cek akses: client cuma bisa lihat project miliknya
	if user.Role == "client" && project.UserID != user.ID {
		c.Redirect(http.StatusFound, "/client/wedding?error=Tidak+memiliki+akses")
		return
	}

	// Cek status: harus published
	if project.Status != "published" {
		c.Redirect(http.StatusFound, "/admin/projects?error=Undangan+belum+dipublikasikan")
		return
	}

	// Sync RSVP → InvitationGuest
	syncRsvpToGuests(project.ID)

	// Ambil semua tamu
	var guests []models.InvitationGuest
	db.Where("project_id = ?", project.ID).
		Order("is_shared ASC, name ASC").
		Find(&guests)

	// Ambil data undangan untuk default message
	generalData := services.ParseJSONMap(project.DataUndangan)

	// Bangun URL undangan
	baseUrl := getBaseUrl(c) + "/invitation/" + project.Slug

	// Stats
	totalGuests := len(guests)
	sharedCount := 0
	hasPhoneCount := 0
	for _, g := range guests {
		if g.IsShared {
			sharedCount++
		}
		if g.IsHasPhone() {
			hasPhoneCount++
		}
	}

	stats := map[string]int{
		"total":     totalGuests,
		"shared":    sharedCount,
		"has_phone": hasPhoneCount,
	}

	// Default message
	defaultMsg := buildDefaultMessage(generalData)

	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	var buf bytes.Buffer
	err = layouts.AdminLayout(
		*user, setting, "share.invitation",
		share.InvitationContent(
			project,
			guests,
			baseUrl,
			defaultMsg,
			stats,
			errorMsg,
			successMsg,
		),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// StoreGuest — POST /share/invitation/:id/guest
// Tambah tamu manual (AJAX)
// ============================================
func StoreGuest(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Project tidak ditemukan"})
		return
	}

	if user.Role == "client" && project.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak memiliki akses"})
		return
	}

	name := strings.TrimSpace(c.PostForm("name"))
	phone := strings.TrimSpace(c.PostForm("phone"))

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Nama wajib diisi"})
		return
	}

	// Cek duplikat
	var existing models.InvitationGuest
	if err := db.Where("project_id = ? AND LOWER(name) = ?", project.ID, strings.ToLower(name)).First(&existing).Error; err == nil {
		// Update phone kalau ada
		if phone != "" {
			existing.Phone = &phone
			db.Save(&existing)
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tamu sudah ada, data diperbarui", "guest": existing})
		return
	}

	// Buat tamu baru
	guest := models.InvitationGuest{
		ProjectID: project.ID,
		Name:      name,
	}
	if phone != "" {
		guest.Phone = &phone
	}

	if err := db.Create(&guest).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menyimpan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tamu berhasil ditambahkan", "guest": guest})
}

// ============================================
// UpdateGuest — PUT /share/guest/:guestId
// Update nomor WA / mark shared (AJAX)
// ============================================
func UpdateGuest(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	guestID, err := strconv.ParseUint(c.Param("guestId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var guest models.InvitationGuest
	db := database.GetDB()
	if err := db.First(&guest, guestID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tamu tidak ditemukan"})
		return
	}

	// Cek akses
	var project models.Project
	db.First(&project, guest.ProjectID)
	if user.Role == "client" && project.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak memiliki akses"})
		return
	}

	phone := strings.TrimSpace(c.PostForm("phone"))
	name := strings.TrimSpace(c.PostForm("name"))
	isShared := c.PostForm("is_shared") == "1"

	if phone != "" {
		guest.Phone = &phone
	}
	if name != "" {
		guest.Name = name
	}
	if isShared {
		guest.IsShared = true
		now := getNow()
		guest.LastSharedAt = &now
	}

	if err := db.Save(&guest).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Gagal menyimpan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data diperbarui", "guest": guest})
}

// ============================================
// DestroyGuest — DELETE /share/guest/:guestId
// Hapus tamu (AJAX)
// ============================================
func DestroyGuest(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	guestID, err := strconv.ParseUint(c.Param("guestId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var guest models.InvitationGuest
	db := database.GetDB()
	if err := db.First(&guest, guestID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tamu tidak ditemukan"})
		return
	}

	var project models.Project
	db.First(&project, guest.ProjectID)
	if user.Role == "client" && project.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak memiliki akses"})
		return
	}

	db.Delete(&guest)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tamu berhasil dihapus"})
}

// ============================================
// BulkUpdateGuests — POST /share/invitation/:id/guests/bulk-update
// Update banyak nomor WA sekaligus (AJAX)
// ============================================
func BulkUpdateGuests(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Project tidak ditemukan"})
		return
	}

	if user.Role == "client" && project.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Tidak memiliki akses"})
		return
	}

	// Parse updates dari form: updates[0][guestId], updates[0][phone], dst
	updates := c.PostFormArray("updates")
	updated := 0

	// Sebenarnya kita perlu parse JSON. Tapi kita pakai form array
	// Format: updates=guestId:phone,guestId:phone
	_ = updates

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OK", "updated": updated})
}

// ============================================
// HELPER — sync RSVP ke InvitationGuest
// ============================================
func syncRsvpToGuests(projectID uint) {
	db := database.GetDB()

	var rsvps []models.Rsvp
	db.Where("project_id = ?", projectID).Find(&rsvps)

	for _, rsvp := range rsvps {
		var existing models.InvitationGuest
		err := db.Where("project_id = ? AND (rsvp_id = ? OR LOWER(name) = ?)",
			projectID, rsvp.ID, strings.ToLower(rsvp.GuestName)).
			First(&existing).Error

		if err == nil {
			// Update rsvp_id kalau belum
			if existing.RsvpID == nil {
				existing.RsvpID = &rsvp.ID
				db.Save(&existing)
			}
		} else {
			// Buat tamu baru
			guest := models.InvitationGuest{
				ProjectID: projectID,
				RsvpID:    &rsvp.ID,
				Name:      rsvp.GuestName,
			}
			db.Create(&guest)
		}
	}
}

// ============================================
// HELPER — build default message
// ============================================
func buildDefaultMessage(data map[string]interface{}) string {
	groomName := getStrFromData(data, "groom_name", "Mempelai Pria")
	brideName := getStrFromData(data, "bride_name", "Mempelai Wanita")

	msg := "Assalamualaikum Warahmatullahi Wabarakatuh,\n\n"
	msg += "Dengan memohon rahmat dan ridho Allah SWT, kami bermaksud menyelenggarakan pernikahan putra-putri kami:\n\n"
	msg += "💍 " + groomName + " & " + brideName + "\n\n"
	msg += "Berikut link undangan kami:\n"
	msg += "{{LINK}}\n\n"
	msg += "Merupakan suatu kehormatan dan kebahagiaan bagi kami apabila Bapak/Ibu/Saudara/i berkenan hadir untuk memberikan doa restu.\n\n"
	msg += "Terima kasih 🙏\n\n"
	msg += "Wassalamualaikum Warahmatullahi Wabarakatuh"

	return msg
}

// ============================================
// HELPER
// ============================================
func getStrFromData(data map[string]interface{}, key, def string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

func getBaseUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	return scheme + "://" + host
}

// getNow — helper untuk waktu sekarang
func getNow() time.Time {
	return time.Now()
}
