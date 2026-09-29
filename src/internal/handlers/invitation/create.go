package invitation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wedding-invitation-go/internal/database"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/internal/services"
	"wedding-invitation-go/views/invitation"
	"wedding-invitation-go/views/layouts"
)

// ============================================
// SelectClient — GET /invitation/create/select-client
// ============================================
func SelectClient(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	if user.Role == "client" {
		services.WizardSetClientID(c.Writer, c.Request, user.ID)
		c.Redirect(http.StatusFound, "/invitation/create/general")
		return
	}

	var clients []models.User
	db := database.GetDB()
	db.Where("role = ?", "client").Order("name ASC").Find(&clients)

	var buf bytes.Buffer
	err := layouts.AdminLayout(
		*user, setting, "invitation.create",
		invitation.SelectClientContent(clients, c.Query("error")),
	).Render(c.Request.Context(), &buf)

	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// ============================================
// StoreClient — POST /invitation/create/select-client
// ============================================
func StoreClient(c *gin.Context) {
	clientIDStr := c.PostForm("client_id")

	var clientID uint
	for _, ch := range clientIDStr {
		if ch < '0' || ch > '9' {
			c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Client+wajib+dipilih")
			return
		}
		clientID = clientID*10 + uint(ch-'0')
	}

	if clientID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Client+wajib+dipilih")
		return
	}

	db := database.GetDB()
	var client models.User
	if err := db.Where("id = ? AND role = ?", clientID, "client").First(&client).Error; err != nil {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Client+tidak+ditemukan")
		return
	}

	services.WizardSetClientID(c.Writer, c.Request, clientID)
	c.Redirect(http.StatusFound, "/invitation/create/general")
}

// ============================================
// EditStart — GET /invitation/create/edit/:id
// Set session wizard dari project existing,
// lalu redirect ke General (form edit)
// ============================================
func EditStart(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// Parse ID dari URL
	idStr := c.Param("id")
	var id uint
	for _, ch := range idStr {
		if ch < '0' || ch > '9' {
			break
		}
		id = id*10 + uint(ch-'0')
	}

	if id == 0 {
		c.Redirect(http.StatusFound, "/admin/projects?error=ID+tidak+valid")
		return
	}

	// Load project
	var project models.Project
	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/projects?error=Project+tidak+ditemukan")
		return
	}

	// Cek akses: client cuma bisa edit project miliknya
	if user.Role == "client" && project.UserID != user.ID {
		c.Redirect(http.StatusFound, "/admin/projects?error=Tidak+memiliki+akses")
		return
	}

	// ✅ Set session wizard: project_id + client_id
	services.WizardSetProjectID(c.Writer, c.Request, project.ID)
	services.WizardSetClientID(c.Writer, c.Request, project.UserID)

	services.LogSuccess("EditStart", fmt.Sprintf("Edit mode untuk project %d (%s)", project.ID, project.Slug))

	// Redirect ke General (mode edit)
	c.Redirect(http.StatusFound, "/invitation/create/general")
}

// ============================================
// General — GET /invitation/create/general
// Step 1: Form data umum + bank accounts
// ============================================
func General(c *gin.Context) {
	setting := services.GetSiteSetting()
	user := c.MustGet("user").(*models.User)

	clientID := services.WizardGetClientID(c.Request)
	if clientID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis,+silakan+mulai+ulang")
		return
	}

	// Ambil project existing (kalau edit)
	projectID := services.WizardGetProjectID(c.Request)
	var project *models.Project
	var dataUndangan map[string]interface{}

	if projectID > 0 {
		var p models.Project
		db := database.GetDB()
		if err := db.First(&p, projectID).Error; err == nil {
			project = &p
			dataUndangan = parseJSONToMap(p.DataUndangan)
		}
	}

	if dataUndangan == nil {
		dataUndangan = make(map[string]interface{})
	}

	errorMsg := c.Query("error")

	renderWithLayout(c, user, setting,
		invitation.GeneralContent(
			project,
			dataUndangan,
			getBankListAsInterface(),
			errorMsg,
			clientID,
			isClient(user),
		),
	)
}

// ============================================
// GeneralSave — POST /invitation/create/general
// Simpan data umum + bank accounts
// ============================================
func GeneralSave(c *gin.Context) {
	clientID := services.WizardGetClientID(c.Request)
	if clientID == 0 {
		c.Redirect(http.StatusFound, "/invitation/create/select-client?error=Sesi+habis")
		return
	}

	// Validasi input
	groomName := strings.TrimSpace(c.PostForm("groom_name"))
	brideName := strings.TrimSpace(c.PostForm("bride_name"))
	akadDate := strings.TrimSpace(c.PostForm("akad_date"))
	akadTime := strings.TrimSpace(c.PostForm("akad_time"))
	akadVenue := strings.TrimSpace(c.PostForm("akad_venue"))
	resepsiDate := strings.TrimSpace(c.PostForm("resepsi_date"))
	resepsiTime := strings.TrimSpace(c.PostForm("resepsi_time"))
	resepsiVenue := strings.TrimSpace(c.PostForm("resepsi_venue"))
	resepsiLabel := strings.TrimSpace(c.PostForm("resepsi_label"))
	resepsiLabelCustom := strings.TrimSpace(c.PostForm("resepsi_label_custom"))
	showDates := strings.TrimSpace(c.PostForm("show_dates"))
	showVenue := strings.TrimSpace(c.PostForm("show_venue"))
	showBankAccounts := strings.TrimSpace(c.PostForm("show_bank_accounts"))

	if groomName == "" || brideName == "" {
		c.Redirect(http.StatusFound, "/invitation/create/general?error=Nama+mempelai+wajib+diisi")
		return
	}
	if akadDate == "" || akadTime == "" || akadVenue == "" {
		c.Redirect(http.StatusFound, "/invitation/create/general?error=Data+akad+wajib+lengkap")
		return
	}
	if resepsiDate == "" || resepsiTime == "" || resepsiVenue == "" {
		c.Redirect(http.StatusFound, "/invitation/create/general?error=Data+resepsi+wajib+lengkap")
		return
	}

	// Format tanggal Indonesia
	akadDateFormatted := formatDateIndonesia(akadDate)
	resepsiDateFormatted := formatDateIndonesia(resepsiDate)

	// Parse bank accounts dari form
	bankAccounts := parseBankAccountsFromForm(c)

	// Bangun data undangan
	dataUndangan := map[string]interface{}{
		"groom_name": groomName,
		"bride_name": brideName,

		"akad_date":     akadDateFormatted,
		"akad_date_raw": akadDate,
		"akad_time":     akadTime,
		"akad_venue":    akadVenue,

		"resepsi_date":     resepsiDateFormatted,
		"resepsi_date_raw": resepsiDate,
		"resepsi_time":     resepsiTime,
		"resepsi_venue":    resepsiVenue,

		"resepsi_label":        resepsiLabel,
		"resepsi_label_custom": resepsiLabelCustom,

		"show_dates":         showDates,
		"show_venue":         showVenue,
		"show_bank_accounts": showBankAccounts,
		"bank_accounts":      bankAccounts,
	}

	dataJSON, _ := json.Marshal(dataUndangan)

	// Cari project existing
	projectID := services.WizardGetProjectID(c.Request)
	db := database.GetDB()

	if projectID > 0 {
		// Update existing project
		var p models.Project
		if err := db.First(&p, projectID).Error; err == nil {
			p.DataUndangan = string(dataJSON)
			p.Title = groomName + " & " + brideName
			if err := db.Save(&p).Error; err != nil {
				fmt.Printf("❌ Gagal update project: %v\n", err)
				errorMsg := strings.ReplaceAll(err.Error(), " ", "+")
				c.Redirect(http.StatusFound, "/invitation/create/general?error="+errorMsg)
				return
			}

			// ✅ FIX: redirect ke step berikutnya & STOP (jangan lanjut create baru)
			services.LogSuccess("GeneralSave", "Project diupdate: "+p.Slug)
			c.Redirect(http.StatusFound, "/invitation/create/template")
			return
		}
	}

	// ↓↓↓ HANYA JALAN KALAU projectID == 0 (create baru) ↓↓↓
	slug := generateSlug(groomName + "-" + brideName)

	// Pastikan slug unik
	var existing models.Project
	if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		slug = slug + "-" + fmt.Sprintf("%d", time.Now().Unix()%10000)
	}

	// Ambil template default (nanti bisa diganti di step 2)
	var defaultTemplate models.Template
	db.Where("is_active = ?", true).Order("`order` ASC").First(&defaultTemplate)

	project := models.Project{
		UserID:               clientID,
		TemplateID:           defaultTemplate.ID,
		Title:                groomName + " & " + brideName,
		Slug:                 slug,
		DataUndangan:         string(dataJSON),
		TemplateSpecificData: "{}", // Default JSON kosong
		Status:               "draft",
	}

	if err := db.Create(&project).Error; err != nil {
		services.LogError("GeneralSave.CreateProject", err, map[string]interface{}{
			"groom_name": groomName,
			"bride_name": brideName,
			"client_id":  clientID,
			"slug":       slug,
			"template":   defaultTemplate.ID,
		})
		c.Redirect(http.StatusFound, "/invitation/create/general?error=Gagal+membuat+project")
		return
	}

	services.LogSuccess("GeneralSave", "Project dibuat: "+project.Slug)

	// Simpan project_id di session
	services.WizardSetProjectID(c.Writer, c.Request, project.ID)

	c.Redirect(http.StatusFound, "/invitation/create/template")
}

// ============================================
// HELPER — parse bank accounts dari form
// ============================================
func parseBankAccountsFromForm(c *gin.Context) []map[string]interface{} {
	var accounts []map[string]interface{}

	// Bank accounts dikirim sebagai: bank_accounts[0][type], bank_accounts[0][bank_code], dst
	// Kita pakai c.PostForm untuk akses dinamis
	for i := 0; i < 10; i++ {
		prefix := fmt.Sprintf("bank_accounts[%d]", i)
		accNumber := strings.TrimSpace(c.PostForm(prefix + "[account_number]"))

		if accNumber == "" {
			continue
		}

		account := map[string]interface{}{
			"id":             strings.TrimSpace(c.PostForm(prefix + "[id]")),
			"type":           strings.TrimSpace(c.PostForm(prefix + "[type]")),
			"bank_code":      strings.TrimSpace(c.PostForm(prefix + "[bank_code]")),
			"bank_name":      strings.TrimSpace(c.PostForm(prefix + "[bank_name]")),
			"account_number": accNumber,
			"account_name":   strings.TrimSpace(c.PostForm(prefix + "[account_name]")),
			"icon_type":      strings.TrimSpace(c.PostForm(prefix + "[icon_type]")),
			"icon_path":      strings.TrimSpace(c.PostForm(prefix + "[existing_icon_path]")),
		}

		if account["id"] == "" {
			account["id"] = fmt.Sprintf("bank_%d_%d", time.Now().Unix(), i)
		}

		accounts = append(accounts, account)
	}

	return accounts
}

// ============================================
// HELPER — format tanggal Indonesia
// ============================================
func formatDateIndonesia(dateStr string) string {
	if dateStr == "" {
		return "-"
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	days := map[time.Weekday]string{
		time.Sunday: "Minggu", time.Monday: "Senin", time.Tuesday: "Selasa",
		time.Wednesday: "Rabu", time.Thursday: "Kamis", time.Friday: "Jumat",
		time.Saturday: "Sabtu",
	}

	months := map[time.Month]string{
		time.January: "Januari", time.February: "Februari", time.March: "Maret",
		time.April: "April", time.May: "Mei", time.June: "Juni",
		time.July: "Juli", time.August: "Agustus", time.September: "September",
		time.October: "Oktober", time.November: "November", time.December: "Desember",
	}

	return fmt.Sprintf("%s, %d %s %d",
		days[t.Weekday()], t.Day(), months[t.Month()], t.Year())
}

// ============================================
// HELPER — generate slug
// ============================================
func generateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "&", "dan")
	result := ""
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result += string(ch)
		}
	}
	// Trim multiple dashes
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}

// ============================================
// HELPER — parse JSON ke map
// ============================================
func parseJSONToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	if jsonStr == "" {
		return result
	}
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

// ============================================
// HELPER — bank list untuk view
// ============================================
func getBankListAsInterface() []map[string]string {
	var result []map[string]string
	for _, b := range models.GetBankList() {
		result = append(result, map[string]string{
			"code": b.Code,
			"name": b.Name,
			"type": b.Type,
			"full": b.Full,
		})
	}
	return result
}
