📘 PANDUAN BIKIN TEMA BARU — Wedding Invitation Go
🎯 Tujuan

Panduan ini untuk developer yang mau bikin tema undangan baru di project Wedding Invitation Go. Setelah baca ini, kamu bisa bikin tema baru yang:

    Tampil di homepage (pilihan template)

    Tampil di wizard Step 3 (pilih template)

    Punya preview /preview/template/:slug

    Berfungsi penuh (RSVP, music, footer, dll)

📐 Arsitektur Project (Ringkas)
text

wedding-invitation-go/src/
├── internal/
│   ├── invitation/                  ← PACKAGE BARU (data + helper global)
│   │   ├── data.go                  ← struct TemplateData
│   │   ├── loader.go                ← LoadData()
│   │   └── helpers.go               ← method & fungsi global
│   │
│   ├── handlers/
│   │   ├── invitation/show.go       ← render undangan publik
│   │   └── preview_template.go      ← render preview template
│   │
│   ├── models/
│   │   ├── project.go
│   │   ├── template.go
│   │   └── bank_account.go
│   │
│   ├── services/
│   │   ├── sitesetting.go
│   │   ├── template_data.go
│   │   └── rsvp.go
│   │
│   └── database/
│       └── template_seed.go         ← SEED TEMPLATE + FIELD LIBRARY
│
├── views/
│   ├── invitations/                 ← TEMPLATE UNDANGAN
│   │   ├── shared/                  ← KOMPONEN SHARED
│   │   │   ├── rsvp.templ           ← form RSVP + list
│   │   │   ├── music_player.templ   ← music player
│   │   │   └── web_footer.templ     ← web footer
│   │   │
│   │   ├── rustic_wood/index.templ  ← tema contoh
│   │   ├── muslim_elegan/index.templ
│   │   ├── elegant_gold/index.templ
│   │   ├── modern_minimalist/index.templ
│   │   ├── traditional_java/index.templ
│   │   └── render.go                ← dispatcher (WAJIB DIUPDATE)
│   │
│   └── invitation/                  ← WIZARD (JANGAN DISENTUH)
│       └── ...
│
└── public/storage/                  ← file upload & placeholder
    └── defaults/                    ← foto placeholder

🔄 Alur Data Undangan
text

1. DB: projects (slug, data_undangan JSON, template_specific_data JSON, template_id)
        templates (folder, fields_schema JSON)

2. Handler ShowInvitation:
   - Query project by slug
   - invitation.LoadData(project, guestName) → *TemplateData
   - RenderTemplate(buf, r, project, data, guestName)

3. render.go dispatcher:
   - Cek project.Template.Folder → panggil tema sesuai

4. Tema:
   - Terima *TemplateData
   - Akses field: data.GroomName, data.AkadDate, dll
   - Panggil shared: @shared.RSVPForm(data, ...), @shared.MusicPlayer(), @shared.WebFooter(data, ...)

📋 Struct TemplateData — Field yang Tersedia

Semua field diakses via data.FieldName. Contoh: data.GroomName, data.AkadDate.
Field Umum (dari data_undangan):
Field	Tipe	Keterangan
GroomName	string	Nama mempelai pria
BrideName	string	Nama mempelai wanita
GroomPhoto	string	Path foto pria (bisa kosong)
BridePhoto	string	Path foto wanita
HeroImage	string	Foto hero/cover
FatherGroom, MotherGroom	string	Nama orang tua pria
FatherBride, MotherBride	string	Nama orang tua wanita
FatherGroomPhoto, MotherGroomPhoto	string	Foto ortu pria
FatherBridePhoto, MotherBridePhoto	string	Foto ortu wanita
GroomInstagram, BrideInstagram	string	IG mempelai
EventInstagram	string	IG acara
Whatsapp	string	Nomor WA
Youtube	string	Link YouTube
AkadDate, AkadTime, AkadVenue	string	Info akad
ResepsiDate, ResepsiTime, ResepsiVenue	string	Info resepsi
ResepsiLabel	string	"resepsi", "walimatul_ursy", "walimah", "custom"
ResepsiLabelCustom	string	Kalau custom
MapsURLAkad, MapsURLResepsi	string	Link Google Maps
ShowDates	string	"both", "akad", "resepsi", "none"
ShowVenue	string	"both", "akad", "resepsi", "none"
ShowBankAccounts	string	"yes", "no"
GroomFamilyOrigin, BrideFamilyOrigin	string	Asal keluarga (untuk adat)
KembarMayang	string	"yes", "no"
Field Parsed (dari JSON, sudah jadi struct/slice):
Field	Tipe	Keterangan
BankAccounts	[]BankAccount	List rekening (sudah parsed)
Gallery	[]string	List URL foto galeri
LoveStories	[]LoveStory	Cerita cinta (Title, Desc)
SiteConfig	models.Config	Setting situs (favicon, sosmed, dll)
ExistingRsvp	*ExistingRsvp	RSVP tamu (kalau ada)
Specific	map[string]interface{}	Field dinamis dari template_specific_data
Project	models.Project	Project data (ID, Slug, dll)
GuestName	string	Nama tamu (dari ?to=)
IsNamedGuest	bool	true kalau tamu punya nama
MusicURL	string	URL musik (sudah full path)
Method Bantuan (sudah tersedia):
go

// Foto dengan fallback default
data.GroomPhotoOrDefault() string        // fallback ke /storage/defaults/groom.svg
data.BridePhotoOrDefault() string
data.FatherGroomPhotoOrDefault() string
data.MotherGroomPhotoOrDefault() string
data.FatherBridePhotoOrDefault() string
data.MotherBridePhotoOrDefault() string
data.HeroImageOr(fallback string) string // HeroImage atau fallback
data.HasHeroImage() bool

// Resepsi label
data.ResepsiLabelDisplay() string  // "Walimatul Ursy", "Resepsi", dll (human-readable)
data.ResepsiArabic() string        // "وَلِيمَةُ الْعُرْسِ" (kalau islami)
data.IsIslamicResepsi() bool       // true kalau "walimatul_ursy" atau "walimah"

// Orang tua
data.ParentsGroom() string  // "Bapak & Ibu" kalau kosong
data.ParentsBride() string

// Tampilan
data.ShouldShowBankAccounts() bool  // true kalau ShowBankAccounts == "yes"
data.ShouldShowAkadDate() bool
data.ShouldShowResepsiDate() bool
data.ShouldShowBothDates() bool
data.ShouldShowAkadVenue() bool
data.ShouldShowResepsiVenue() bool
data.ShouldShowBothVenues() bool

// RSVP
data.HasExistingRsvp() bool
data.ExistingRsvpField(key, def string) string  // key: "attendance" / "message"
data.ExistingRsvpGuests() int                    // fallback 1

Fungsi Statis (untuk bank & format):
go

// Import: "wedding-invitation-go/internal/invitation"
invitation.GetBankTypeBadge(acc) string        // "BANK", "E-WALLET", "QRIS", "LAINNYA"
invitation.GetAccountBankName(acc) string      // "BCA", "Mandiri", dll
invitation.GetAccountInitial(acc) string       // 2 huruf awal
invitation.GetAccountNumber(acc) string        // nomor atau "-"
invitation.GetAccountName(acc) string          // nama atau "-"
invitation.HasCustomIcon(acc) bool
invitation.GetAccountIconURL(acc) string       // full path icon

// Format
invitation.FormatUint(n uint) string
invitation.FormatInt(n int) string
invitation.CurrentYear() int
invitation.GuestGreeting(isNamed bool) string  // "Kepada" / "Salam Hangat untuk"

📦 Komponen Shared — WAJIB DIPAKAI
1. @shared.RSVPForm(data, text)

Form RSVP + list. Parameter text:

    shared.RSVPTextIndonesia — "Nama Anda", "Kehadiran", "Hadir", dll

    shared.RSVPTextJawa — "Nama", "Rawuh", "Mboten Rawuh", dll

Atau buat custom shared.RSVPText{...}.
2. @shared.RSVPListWrap()

Wrapper untuk list konfirmasi di bawah form. Wajib dipanggil setelah RSVPForm.
3. @shared.MusicPlayer()

Music player tombol bulat di kanan bawah.
4. @shared.WebFooter(data, text)

Footer web (brand, kontak, developer). Parameter text:

    shared.WebFooterTextDefault — "Kontak Kami", "Dikembangkan Oleh"

    shared.WebFooterTextShort — "Kontak", "Developer"

🎨 CSS — Kontrak Class

Class CSS WAJIB ada di tema baru (shared mengandalkan ini):
RSVP Form:
css

.rsvp-form { /* style form */ }
.form-group { /* field group */ }
.form-group label { /* label */ }
.form-group input, .form-group select, .form-group textarea { /* input */ }
.form-group .required { /* tanda * */ }
.btn-submit { /* tombol submit */ }
.rsvp-warning { /* warning sudah RSVP */ }
.rsvp-warning-title { /* judul warning */ }
.rsvp-warning-title i { /* icon */ }
.rsvp-warning p { /* teks warning */ }
.rsvp-list-wrap { /* wrapper list */ }
.rsvp-list-head { /* head list */ }
.rsvp-list-head h4 { /* judul list */ }
.rsvp-count { /* count badge */ }
.rsvp-empty { /* list kosong */ }
.rsvp-item { /* 1 item RSVP */ }
.rsvp-item .row-main { /* baris utama */ }
.rsvp-item .guest-name { /* nama tamu */ }
.rsvp-item .status-badge { /* badge status */ }
.rsvp-item .status-badge.hadir { /* hijau */ }
.rsvp-item .status-badge.tidak_hadir { /* merah */ }
.rsvp-item .status-badge.ragu { /* kuning */ }
.rsvp-item .row-detail { /* baris detail */ }
.rsvp-item .meta-left { /* meta kiri */ }
.rsvp-item .meta-left .time { /* waktu */ }
.rsvp-item .meta-left .guests { /* jumlah tamu */ }
.rsvp-item .message { /* pesan */ }

Music Player:
css

.music-player { position: fixed; bottom: 24px; right: 24px; z-index: 1000; }
.music-btn { /* tombol bulat */ }
.music-btn.playing { /* animasi saat playing */ }

Web Footer:
css

.web-footer { /* container */ }
.web-footer-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; }
.web-footer-brand { /* brand */ }
.web-footer-brand .wedding { /* warna */ }
.web-footer-brand .saas { /* warna */ }
.web-footer-desc { /* deskripsi */ }
.web-footer-link { /* link */ }
.web-footer-socials { /* social container */ }
.web-social-link { /* social link */ }
.web-footer-heading { /* heading */ }
.web-footer-list { /* list */ }
.web-footer-list li { /* item */ }
.web-footer-list li i { /* icon */ }
.web-footer-version { /* versi */ }
.web-footer-divider { /* divider */ }
.web-footer-text { /* copyright */ }
.web-footer-icon { color: var(--warna-tema); }

Toast (dibuat via JS, tapi CSS perlu):
css

.toast-container { position: fixed; top: 20px; right: 20px; z-index: 9999; }
.toast { /* style */ }
.toast.success i { color: #10b981; }
.toast.error i { color: #ef4444; }
.toast.info i { color: var(--warna-tema); }

Overlay (khusus tema, tidak shared):
css

#startOverlay { /* overlay */ }
.overlay-inner { /* inner */ }
.btn-start { /* tombol buka */ }

📝 JavaScript — Kontrak Fungsi

JS WAJIB ada di tema baru (dipanggil shared & HTML):
Wajib karena shared panggil:
js

window.startMusic = function() { ... }    // dipanggil dari overlay + music player
window.toggleMusic = function() { ... }   // dipanggil dari music player

Wajib karena HTML tema panggil:
js

window.scrollToSection = function() { ... } // untuk tombol scroll (bisa di-alias)

JS standar yang bisa copy dari tema lain:

    AOS init

    Music player (startMusic, toggleMusic, closeOverlay, initAudio, retry)

    Toast (showToast, escapeHtml)

    Copy bank (copyBankNumberFromEl, copyText, fallbackCopy)

    RSVP load + submit (loadRsvpList)

    Lightbox

    Scroll helper

Cara termudah: copy dari rustic_wood/index.templ, ganti nama class gallery & section ID.
🚀 CARA BIKIN TEMA BARU (Step-by-Step)
Step 1: Buat folder & file
bash

mkdir -p /DATA/AppData/wedding-invitation-go/src/views/invitations/tema_baru
cd /DATA/AppData/wedding-invitation-go/src/views/invitations/tema_baru

Buat 1 file saja: index.templ
Step 2: Struktur file index.templ
templ

package tema_baru

import (
	"context"
	"io"

	templpkg "github.com/a-h/templ"
	"wedding-invitation-go/internal/invitation"
	"wedding-invitation-go/internal/models"
	"wedding-invitation-go/views/invitations/shared"
)

// TemaBaru — fungsi entry point
func TemaBaru(ctx context.Context, w io.Writer, project models.Project, data *invitation.TemplateData, guestName string) error {
	component := temaBaruPage(data)
	return component.Render(ctx, w)
}

templ temaBaruPage(data *invitation.TemplateData) {
	<!DOCTYPE html>
	<html lang="id">
	<head>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>{ data.GroomName } &amp; { data.BrideName } - Tema Baru</title>
		
		<!-- Fonts, CSS libraries (bootstrap-icons, AOS) -->
		<link href="https://fonts.googleapis.com/css2?family=..." rel="stylesheet"/>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css"/>
		<link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet"/>
		
		<style>
			/* CSS RESET */
			* { margin: 0; padding: 0; box-sizing: border-box; }
			body { font-family: 'Inter', sans-serif; background: var(--bg); }
			
			/* CSS VARIABLES */
			:root {
				--warna-tema: #...;
				/* ... */
			}
			
			/* CSS TEMA (contoh: .hero, .couple, dll) */
			/* WAJIB: RSVP, music player, web footer — copy dari tema lain & adjust warna */
			/* ... */
		</style>
	</head>
	<body data-music-url={ data.MusicURL } data-project-id={ invitation.FormatUint(data.Project.ID) } data-project-slug={ data.Project.Slug }>

		<!-- OVERLAY -->
		<div id="startOverlay">
			...desain overlay tema baru...
		</div>

		<!-- MUSIC PLAYER — WAJIB -->
		@shared.MusicPlayer()

		<!-- HERO -->
		<section class="hero">
			...desain hero tema baru...
		</section>

		<!-- COUPLE -->
		<section class="section-couple">
			...desain couple tema baru...
		</section>

		<!-- LOVE STORY (opsional) -->
		if len(data.LoveStories) > 0 {
			<section>
				for _, story := range data.LoveStories {
					<div>{ story.Title } — { story.Desc }</div>
				}
			</section>
		}

		<!-- EVENT -->
		<section>
			<div>{ data.AkadDate } — { data.AkadVenue }</div>
			<div>{ data.ResepsiDate } — { data.ResepsiVenue }</div>
		</section>

		<!-- GALLERY (opsional) -->
		if len(data.Gallery) > 0 {
			<section>
				for _, img := range data.Gallery {
					<img src={ img }/>
				}
			</section>
		}

		<!-- RSVP — WAJIB -->
		<section class="section-rsvp">
			<div class="container">
				<h2>Konfirmasi Kehadiran</h2>
				@shared.RSVPForm(data, shared.RSVPTextIndonesia)
				@shared.RSVPListWrap()
			</div>
		</section>

		<!-- BANK / AMPLOP DIGITAL (opsional) -->
		if data.ShouldShowBankAccounts() && len(data.BankAccounts) > 0 {
			<section>
				for _, acc := range data.BankAccounts {
					<div>{ invitation.GetAccountBankName(acc) } — { invitation.GetAccountNumber(acc) }</div>
				}
			</section>
		}

		<!-- WEB FOOTER — WAJIB -->
		@shared.WebFooter(data, shared.WebFooterTextDefault)

		<!-- SCRIPTS — copy dari tema lain, adjust -->
		<script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
		<script>
			// AOS init
			// Music player
			// Toast
			// Copy bank
			// RSVP load & submit
			// Lightbox
			// Scroll helper
		</script>
	</body>
	</html>
}

Step 3: Register di render.go

Edit views/invitations/render.go:
go

import (
	// ... import lain
	"wedding-invitation-go/views/invitations/tema_baru"  // ← tambah ini
)

func RenderTemplate(...) error {
	// ...
	switch folder {
	case "rustic_wood": ...
	case "muslim_elegan": ...
	case "elegant_gold": ...
	case "modern_minimalist": ...
	case "traditional_java": ...
	case "tema_baru":  // ← tambah case ini
		return tema_baru.TemaBaru(ctx, w, project, data, guestName)
	default:
		return fmt.Errorf("template %s tidak dikenal", folder)
	}
}

Step 4: Register di DB (template_seed.go)

Buka file: internal/database/template_seed.go

Cari fungsi seedTemplates() — tambah di akhir:
go

// ============================================
// 6. TEMA BARU
// ============================================
temaBaruFields := []string{}
for _, f := range defaultFields {
	temaBaruFields = append(temaBaruFields, f)
	// Tambah field unik kalau perlu:
	// if f == "mother_bride" {
	//     temaBaruFields = append(temaBaruFields, "field_unik_tema_baru")
	// }
}
if err := upsertTemplate(
	"Tema Baru",              // Nama tampil
	"tema-baru",              // Slug (URL-friendly, huruf kecil, tanda hubung)
	"tema_baru",              // Folder (harus sama dengan nama folder di views/invitations/)
	"Deskripsi singkat tema.", // Deskripsi tampil di pilihan template
	buildSchema(temaBaruFields, library),
	6,                        // Order (urutan di pilihan template)
); err != nil {
	return err
}
log.Println("✅ Template: Tema Baru")

Kalau field unik baru, tambahkan dulu di getFieldLibrary():
go

func getFieldLibrary() map[string]FieldDef {
	return map[string]FieldDef{
		// ... field yang sudah ada
		"field_unik_tema_baru": {
			"name":        "field_unik_tema_baru",
			"label":       "Label Field Unik",
			"type":        "text",  // atau "file", "select", "repeater"
			"required":    false,
			"placeholder": "Contoh: ...",
			"group":       "🎨 Khusus Tema Baru",
			"description": "Deskripsi field",
		},
	}
}

Step 5: Build & Restart
bash

cd /DATA/AppData/wedding-invitation-go/src

# Generate templ
docker exec -it wedding-invitation-go-app sh -c "cd /app && templ generate 2>&1 | tail -3"

# Build
docker exec -it wedding-invitation-go-app sh -c "cd /app && go build ./... 2>&1 | head -20"

# Restart
docker restart wedding-invitation-go-app
sleep 10

Step 6: Seed template ke DB

Seed dijalankan otomatis saat startup (lihat main.go). Kalau tidak otomatis, jalankan manual:
bash

docker exec -it wedding-invitation-go-app sh -c "cd /app && go run cmd/seed/main.go"
# atau cara lain sesuai setup kamu

Atau langsung insert via DB:
bash

docker exec wedding-invitation-go-mysql mysql -uwedding_user -pwedding123 wedding_invitation_db -e "
INSERT INTO templates (name, slug, folder, sections, description, fields, fields_schema, is_active, \`order\`, created_at, updated_at)
VALUES (
  'Tema Baru',
  'tema-baru',
  'tema_baru',
  '[\"cover\",\"couple\",\"event\",\"gallery\",\"rsvp\"]',
  'Deskripsi singkat tema.',
  '[]',
  '[]',
  1,
  6,
  NOW(),
  NOW()
);
"

Atau lewat admin panel: /admin/templates → Tambah Template.
Step 7: Test

    Test preview:
    text

    https://wedding.litebox.my.id/preview/template/tema-baru

    Test bikin project baru pakai template ini via wizard:

        Buka /invitation/create/select-client

        Pilih client → Next

        Isi data umum → Next

        Pilih "Tema Baru" → Next

        Isi data spesifik (kalau ada) → Next

        Preview → Publish

    Cek di homepage — pastikan tema baru muncul di pilihan template.

✅ Checklist Tema Baru

    □

    Folder views/invitations/tema_baru/ dibuat
    □

    File index.templ dibuat
    □

    func TemaBaru(ctx, w, project, data, guestName) error didefinisikan
    □

    templ temaBaruPage(data *invitation.TemplateData) didefinisikan
    □

    Import: shared, invitation, templpkg
    □

    <body> punya data-music-url, data-project-id, data-project-slug
    □

    @shared.MusicPlayer() dipanggil
    □

    @shared.RSVPForm(data, shared.RSVPTextIndonesia) dipanggil
    □

    @shared.RSVPListWrap() dipanggil setelah form
    □

    @shared.WebFooter(data, shared.WebFooterTextDefault) dipanggil
    □

    CSS: rsvp-form, form-group, btn-submit, rsvp-warning, rsvp-list-wrap, web-footer-*, music-player, toast-container didefinisikan
    □

    CSS: .web-footer-icon { color: var(--warna-tema); }
    □

    JS: startMusic, toggleMusic, showToast, escapeHtml, copyBankNumberFromEl, loadRsvpList, dll
    □

    render.go diupdate (case baru)
    □

    template_seed.go diupdate (upsertTemplate)
    □

    templ generate sukses
    □

    go build ./... sukses
    □

    Seed template ke DB
    □

    Test preview /preview/template/:slug
    □

    Test wizard create project
    □

    Cek tampil di homepage

🎨 Tips Desain

    CSS Variables — selalu definisikan di :root untuk warna tema:
    css

    :root {
        --primary: #warna-utama;
        --accent: #warna-aksen;
        --bg: #background;
        --text: #teks;
        --text-muted: #teks-lembut;
    }

    Font — beda per tema, pilih yang sesuai karakter:

        Elegan: Playfair Display, Cormorant Garamond

        Modern: Inter, Poppins

        Islami: Amiri, Scheherazade

        Adat: serif klasik

    Ornamen — beda per tema:

        Rustic: ✦, ❦, ⏚

        Islami: ﷲ, ﷻ, ❦

        Minimalis: —, ·, •

        Jawa: ⚜, ❦

    Section unik tema — kalau tema baru punya section khusus (misal "Prosesi Adat"), taruh hanya di tema baru — tidak perlu di-share.

    Test responsive — cek di mobile (max-width 600px) & tablet (600-1024px).

🚫 Yang TIDAK Boleh Diubah

    Jangan sentuh views/invitation/ — itu wizard, bukan template undangan

    Jangan edit file *_templ.go — itu generated, edit .templ saja

    Jangan pakai style={ background-image: ... } — templ escape ', pakai <img> tag

    Jangan pakai {{ var }} di dalam <script> — templ tidak replace {}, pakai data-* attribute

    Jangan tambah helpers.go per tema — semua helper sudah global di internal/invitation/

    Jangan bikin field baru tanpa update getFieldLibrary() — field harus terdaftar

🐛 Debugging Umum
Problem	Solusi
templ generate skip file	Hapus _templ.go, generate ulang
Field tidak tampil	Cek data.FieldName sudah ada di struct? Cek LoadData
Preview 500 error	Cek log: docker logs wedding-invitation-go-app
Template tidak muncul di homepage	Cek DB: SELECT * FROM templates WHERE is_active=1;
Template tidak muncul di wizard	Sama — cek DB
Preview blank	Cek error JS di browser console (F12)
Gambar tidak muncul	Cek path: /storage/... atau URL eksternal
Music tidak play	Cek data-music-url di <body>, cek console
📞 Kontak & Referensi

    Repo: https://github.com/falah-udin/wedding-invitation-go

    Domain: https://wedding.litebox.my.id

    Preview template: /preview/template/:slug

    Contoh tema: lihat views/invitations/rustic_wood/index.templ (paling sederhana)

Selamat berkarya! Kalau ada pertanyaan, tanya ke tim. 🚀
📋 Ringkasan File yang Perlu Dibuat/Diubah

Untuk bikin tema baru:

    views/invitations/tema_baru/index.templ — BARU (~500-1000 baris)

    views/invitations/render.go — UPDATE (1 case)

    internal/database/template_seed.go — UPDATE (1 upsertTemplate + opsional field library)

Total: 1 file baru + 2 file edit. 🎉