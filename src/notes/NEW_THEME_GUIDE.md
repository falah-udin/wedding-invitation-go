📘 PANDUAN BIKIN TEMA BARU — Wedding Invitation Go
Versi: 2.0 (setelah Fase 2 — Shared Components)
Terakhir update: Setelah migrasi RSVP, Music Player, Web Footer ke shared

🎯 Tujuan
Panduan ini untuk developer yang mau bikin tema undangan baru di project Wedding Invitation Go. Setelah baca ini, kamu bisa bikin tema baru yang:

Tampil di homepage (pilihan template)

Tampil di wizard Step 3 (pilih template)

Punya preview /preview/template/:slug

Berfungsi penuh (RSVP, music, footer, dll)

Cukup 1 file index.templ — tidak perlu helper, tidak perlu komponen

📐 Arsitektur Project
text
wedding-invitation-go/src/
├── internal/
│   ├── invitation/                  ← PACKAGE GLOBAL (data + helper)
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
│   │   ├── shared/                  ← KOMPONEN SHARED (JANGAN TULIS MANUAL)
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
│
└── public/storage/                  ← file upload & placeholder
    └── defaults/                    ← foto placeholder (groom.svg, bride.svg, dll)
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

4. Tema (index.templ):
   - Terima *TemplateData
   - Akses field: data.GroomName, data.AkadDate, dll
   - Panggil shared components (WAJIB):
     @shared.MusicPlayer()
     @shared.RSVPForm(data, shared.RSVPTextIndonesia)
     @shared.RSVPListWrap()
     @shared.WebFooter(data, shared.WebFooterTextDefault)
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
Field Parsed (dari JSON):
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
data.HeroImageOr(fallback string) string
data.HasHeroImage() bool

// Resepsi label
data.ResepsiLabelDisplay() string  // "Walimatul Ursy", "Resepsi", dll
data.ResepsiArabic() string        // "وَلِيمَةُ الْعُرْسِ"
data.IsIslamicResepsi() bool       // true kalau islami

// Orang tua
data.ParentsGroom() string  // "Bapak & Ibu" kalau kosong
data.ParentsBride() string

// Tampilan
data.ShouldShowBankAccounts() bool
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
⚠️ ATURAN UTAMA
JANGAN tulis manual:

❌ Form RSVP (<form class="rsvp-form">...</form>)

❌ List wrapper (<div id="rsvpList">...</div>)

❌ Music player (<div class="music-player">...</div>)

❌ Web footer (<footer class="web-footer">...</footer>)

GUNAKAN shared:

✅ @shared.RSVPForm(data, text) + @shared.RSVPListWrap()

✅ @shared.MusicPlayer()

✅ @shared.WebFooter(data, text)

Kenapa? Supaya bug/update cukup di 1 tempat (bukan 5 tema).

1. @shared.RSVPForm(data, text) + @shared.RSVPListWrap()
Cara pakai:

templ
<section class="section-rsvp">
    <div class="container">
        <h2>Konfirmasi Kehadiran</h2>
        @shared.RSVPForm(data, shared.RSVPTextIndonesia)
        @shared.RSVPListWrap()
    </div>
</section>
Parameter text:

shared.RSVPTextIndonesia — "Nama Anda", "Kehadiran", "Hadir", dll

shared.RSVPTextJawa — "Nama", "Rawuh", "Mboten Rawuh", dll

Atau custom: shared.RSVPText{LabelName: "...", ...}

⚠️ WAJIB panggil @shared.RSVPListWrap() SETELAH @shared.RSVPForm() — karena list konfirmasi di bawah form.

2. @shared.MusicPlayer()
Cara pakai (di bawah <body>, setelah overlay):

templ
@shared.MusicPlayer()
Parameter: tidak ada. Cukup panggil.

Prasyarat: <body> harus punya data-music-url, data-project-id, data-project-slug:

templ
<body data-music-url={ data.MusicURL } data-project-id={ invitation.FormatUint(data.Project.ID) } data-project-slug={ data.Project.Slug }>
3. @shared.WebFooter(data, text)
Cara pakai (sebelum </body>):

templ
@shared.WebFooter(data, shared.WebFooterTextDefault)
Parameter text:

shared.WebFooterTextDefault — "Kontak Kami", "Dikembangkan Oleh"

shared.WebFooterTextShort — "Kontak", "Developer"

🎨 CSS — Kontrak Class
Class CSS WAJIB ada di tema baru (shared mengandalkan ini).

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
Copy dari tema lain (contoh: rustic_wood) & adjust warna.

Music Player:
css
.music-player { position: fixed; bottom: 24px; right: 24px; z-index: 1000; }
.music-btn { /* tombol bulat */ }
.music-btn.playing { /* animasi saat playing */ }
Web Footer (⚠️ PERHATIAN KHUSUS):
css
.web-footer { /* container */ }

/* WAJIB — biar 3 kolom */
.web-footer-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
}

/* WAJIB — biar container cukup lebar untuk 3 kolom */
.web-footer .container {
    max-width: 1280px !important;
}

/* WAJIB — biar logo + teks satu baris */
.web-footer-brand img {
    display: inline-block !important;
    vertical-align: middle;
    margin-right: 8px;
}

/* WAJIB — warna icon */
.web-footer-icon { color: var(--warna-tema); }

/* Sisanya */
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
⚠️ PENTING: 3 baris !important di atas WAJIB — kalau tidak:

.web-footer .container tanpa !important → footer tetap 1 kolom (karena .container global lebih dulu)

.web-footer-brand img tanpa !important → logo & teks atas-bawah (karena img { display: block } global)

Kontrak --warna-tema: Ganti dengan CSS var tema kamu, misal var(--rustic-gold), var(--gold), var(--accent-light).

Toast (dibuat via JS, CSS perlu):
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
⚠️ WAJIB ada di tema baru:
js
// Shared panggil:
window.startMusic = function() { ... }
window.toggleMusic = function() { ... }

// HTML tema panggil (nama fungsi bisa di-alias):
window.scrollToSection = function() { ... }
window.scrollToCouple = function() { ... }  // atau nama lain
window.scrollToOpening = function() { ... }
JS standar yang bisa copy dari tema lain:
AOS init — AOS.init({...})

Music player — startMusic, toggleMusic, closeOverlay, initAudio, retry

Toast — showToast, escapeHtml

Copy bank — copyBankNumberFromEl, copyText, fallbackCopy

RSVP — loadRsvpList, submit handler

Lightbox — gallery click handler

Scroll helper — scrollToSection

Cara termudah: Copy dari rustic_wood/index.templ, ganti nama class gallery & section ID.

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

// TemaBaru — entry point
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

        <!-- Fonts -->
        <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;1,400&family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet"/>
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css"/>
        <link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet"/>

        <style>
            /* CSS RESET */
            * { margin: 0; padding: 0; box-sizing: border-box; }
            html { scroll-behavior: smooth; }
            body {
                font-family: 'Inter', sans-serif;
                background: var(--bg);
                color: var(--text);
                line-height: 1.7;
                overflow-x: hidden;
            }
            img { max-width: 100%; height: auto; }

            /* CSS VARIABLES */
            :root {
                --primary: #warna-utama;
                --accent: #warna-aksen;
                --bg: #background;
                --text: #teks;
                --text-muted: #teks-lembut;
            }

            /* CSS TEMA (hero, couple, event, gallery, dll) */
            /* ... */

            /* ============================================ */
            /* RSVP FORM (WAJIB — copy dari tema lain) */
            /* ============================================ */
            .rsvp-form { /* ... */ }
            .form-group { /* ... */ }
            /* dst. (lihat kontrak CSS di atas) */

            /* ============================================ */
            /* MUSIC PLAYER (WAJIB) */
            /* ============================================ */
            .music-player { position: fixed; bottom: 24px; right: 24px; z-index: 1000; }
            .music-btn { /* ... */ }
            .music-btn.playing { /* ... */ }

            /* ============================================ */
            /* WEB FOOTER (WAJIB) */
            /* ============================================ */
            .web-footer { /* ... */ }
            .web-footer-row {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
                gap: 2rem;
            }
            .web-footer .container {
                max-width: 1280px !important;
            }
            .web-footer-brand img {
                display: inline-block !important;
                vertical-align: middle;
                margin-right: 8px;
            }
            .web-footer-icon { color: var(--accent); }
            /* ... dst */

            /* ============================================ */
            /* TOAST */
            /* ============================================ */
            .toast-container { position: fixed; top: 20px; right: 20px; z-index: 9999; }
            .toast { /* ... */ }

            /* ============================================ */
            /* OVERLAY (khusus tema) */
            /* ============================================ */
            #startOverlay { /* ... */ }
        </style>
    </head>
    <body data-music-url={ data.MusicURL } data-project-id={ invitation.FormatUint(data.Project.ID) } data-project-slug={ data.Project.Slug }>

        <!-- OVERLAY (khusus tema) -->
        <div id="startOverlay">
            <!-- desain overlay -->
            <button class="btn-start" onclick={ templpkg.JSFuncCall("startMusic") }>
                Buka Undangan
            </button>
        </div>

        <!-- MUSIC PLAYER — WAJIB -->
        @shared.MusicPlayer()

        <!-- HERO (khusus tema) -->
        <section class="hero">
            <!-- desain hero -->
        </section>

        <!-- COUPLE (khusus tema) -->
        <section class="section-couple">
            <!-- desain couple -->
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

        <!-- BANK (opsional) -->
        if data.ShouldShowBankAccounts() && len(data.BankAccounts) > 0 {
            <section>
                for _, acc := range data.BankAccounts {
                    <div>{ invitation.GetAccountBankName(acc) } — { invitation.GetAccountNumber(acc) }</div>
                }
            </section>
        }

        <!-- WEB FOOTER — WAJIB -->
        @shared.WebFooter(data, shared.WebFooterTextDefault)

        <!-- SCRIPTS -->
        <script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>
        <script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
        <script>
            // AOS init
            AOS.init({ duration: 800, once: true });

            // Music player
            var musicUrl = document.body.dataset.musicUrl || '';
            var projectID = document.body.dataset.projectId || '';
            var projectSlug = document.body.dataset.projectSlug || '';
            window.startMusic = function() { /* ... */ };
            window.toggleMusic = function() { /* ... */ };

            // Toast
            function showToast(type, message) { /* ... */ }
            function escapeHtml(text) { /* ... */ }

            // Copy bank
            window.copyBankNumberFromEl = function(el) { /* ... */ };

            // RSVP load & submit
            function loadRsvpList() { /* ... */ }

            // Lightbox (kalau ada galeri)
            // ...

            // Scroll helper
            window.scrollToSection = function() { /* ... */ };
        </script>
    </body>
    </html>
}
⚠️ Cara tercepat: copy rustic_wood/index.templ, rename package & function, adjust CSS warna & ornamen.

Step 3: Register di render.go
Edit views/invitations/render.go:

go
import (
    // ... import lain
    "wedding-invitation-go/views/invitations/tema_baru"
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
Buka internal/database/template_seed.go, cari fungsi seedTemplates(), tambah di akhir:

go
// ============================================
// 6. TEMA BARU
// ============================================
temaBaruFields := []string{}
for _, f := range defaultFields {
    temaBaruFields = append(temaBaruFields, f)
}
if err := upsertTemplate(
    "Tema Baru",               // Nama tampil
    "tema-baru",               // Slug (URL-friendly)
    "tema_baru",               // Folder (harus sama dengan nama folder)
    "Deskripsi singkat tema.", // Deskripsi
    buildSchema(temaBaruFields, library),
    6,                         // Order
); err != nil {
    return err
}
log.Println("✅ Template: Tema Baru")
Kalau field unik baru, tambahkan dulu di getFieldLibrary().

Step 5: Build & Restart
bash
cd /DATA/AppData/wedding-invitation-go/src

docker exec -it wedding-invitation-go-app sh -c "cd /app && templ generate 2>&1 | tail -3"
docker exec -it wedding-invitation-go-app sh -c "cd /app && go build ./... 2>&1 | head -20"
docker restart wedding-invitation-go-app
sleep 10
Step 6: Seed Template ke DB
Seed dijalankan otomatis saat startup. Kalau tidak otomatis:

bash
docker exec -it wedding-invitation-go-app sh -c "cd /app && go run cmd/seed/main.go"
Atau via admin panel: /admin/templates → Tambah Template.

Step 7: Test
Preview:

text
https://wedding.litebox.my.id/preview/template/tema-baru
Wizard:

Buka /invitation/create/select-client

Pilih client → Next

Isi data umum → Next

Pilih "Tema Baru" → Next

Isi data spesifik → Next

Preview → Publish

Homepage: cek tema baru muncul di pilihan template.

✅ Checklist Tema Baru
□ Folder views/invitations/tema_baru/ dibuat
□ File index.templ dibuat
□ func TemaBaru(ctx, w, project, data, guestName) error didefinisikan
□ templ temaBaruPage(data *invitation.TemplateData) didefinisikan
□ Import: shared, invitation, templpkg
□ <body> punya data-music-url, data-project-id, data-project-slug
Shared components:

□ @shared.MusicPlayer() dipanggil
□ @shared.RSVPForm(data, shared.RSVPTextIndonesia) dipanggil
□ @shared.RSVPListWrap() dipanggil SETELAH RSVPForm
□ @shared.WebFooter(data, shared.WebFooterTextDefault) dipanggil
CSS WAJIB:

□ .rsvp-form, .form-group, .btn-submit, .rsvp-warning, .rsvp-list-wrap
□ .music-player, .music-btn
□ .web-footer, .web-footer-row
□ .web-footer .container { max-width: 1280px !important }
□ .web-footer-brand img { display: inline-block !important; vertical-align: middle }
□ .web-footer-icon { color: var(--warna-tema) }
□ .toast-container, .toast
JS WAJIB:

□ window.startMusic, window.toggleMusic
□ window.scrollToSection (atau alias)
□ showToast, escapeHtml
□ copyBankNumberFromEl, copyText, fallbackCopy
□ loadRsvpList + submit handler
Registrasi:

□ render.go diupdate (case baru)
□ template_seed.go diupdate (upsertTemplate)
Build & Test:

□ templ generate sukses
□ go build ./... sukses
□ Seed template ke DB
□ Test preview /preview/template/:slug
□ Test wizard create project
□ Cek tampil di homepage
🎨 Tips Desain
CSS Variables — selalu definisikan di :root:

css
:root {
    --primary: #warna-utama;
    --accent: #warna-aksen;
    --bg: #background;
    --text: #teks;
    --text-muted: #teks-lembut;
}
Font — beda per tema:

Elegan: Playfair Display, Cormorant Garamond

Modern: Inter, Poppins

Islami: Amiri, Scheherazade

Adat: serif klasik

Ornamen — beda per tema:

Rustic: ✦, ❦, ⏚

Islami: ﷲ, ﷻ, ❦

Minimalis: —, ·, •

Jawa: ⚜, ❦

Test responsive — cek di mobile (max-width 600px) & tablet (600-1024px).

Copy dari tema lain — cara tercepat: copy rustic_wood/index.templ, ganti nama package & function, adjust CSS warna & ornamen.

🚫 Yang TIDAK Boleh Diubah
❌ Jangan sentuh views/invitation/ — itu wizard, bukan template undangan

❌ Jangan edit file *_templ.go — itu generated, edit .templ saja

❌ Jangan pakai style={ background-image: ... } — templ escape ', pakai <img> tag

❌ Jangan pakai {{ var }} di dalam <script> — templ tidak replace {}, pakai data-* attribute

❌ Jangan tambah helpers.go per tema — semua helper global di internal/invitation/

❌ Jangan tulis form RSVP / music player / web footer manual — pakai shared

❌ Jangan lupa panggil @shared.RSVPListWrap() setelah @shared.RSVPForm()

❌ Jangan bikin field baru tanpa update getFieldLibrary() — field harus terdaftar

🐛 Debugging Umum
Problem	Solusi
templ generate skip file	Hapus _templ.go, generate ulang
Field tidak tampil	Cek data.FieldName sudah ada di struct? Cek LoadData
Preview 500 error	Cek log: docker logs wedding-invitation-go-app
Template tidak muncul di homepage	Cek DB: SELECT * FROM templates WHERE is_active=1;
Preview blank	Cek error JS di browser console (F12)
Gambar tidak muncul	Cek path: /storage/... atau URL eksternal
Music tidak play	Cek data-music-url di <body>, cek console
Footer 1 kolom, bukan 3	Cek .web-footer .container { max-width: 1280px !important }
Logo + teks footer atas-bawah	Cek .web-footer-brand img { display: inline-block !important }
RSVP form hilang	Cek @shared.RSVPForm + @shared.RSVPListWrap dipanggil
Music player tidak muncul	Cek @shared.MusicPlayer() dipanggil
📞 Kontak & Referensi
Repo: https://github.com/falah-udin/wedding-invitation-go

Domain: https://wedding.litebox.my.id

Preview template: /preview/template/:slug

Contoh tema: views/invitations/rustic_wood/index.templ (paling sederhana)

📋 Ringkasan File yang Perlu Dibuat/Diubah
Untuk bikin tema baru:

views/invitations/tema_baru/index.templ — BARU (~500-1000 baris)

views/invitations/render.go — UPDATE (1 case)

internal/database/template_seed.go — UPDATE (1 upsertTemplate)

Total: 1 file baru + 2 file edit. 🎉

Selamat berkarya! Kalau ada pertanyaan, tanya ke tim. 🚀

📌 Changelog
v2.0 (setelah Fase 2):

Tambah section "Komponen Shared — WAJIB DIPAKAI"

Tambah CSS wajib untuk web footer (.web-footer .container !important, .web-footer-brand img !important)

Tambah JS kontrak (startMusic, toggleMusic, scrollToSection)

Tambah 3 baris debugging footer (1 kolom, logo atas-bawah, dll)

Update "Yang TIDAK Boleh Diubah" — larangan tulis manual komponen shared

Update checklist dengan item shared components

v1.0 (setelah Fase 1):

Versi awal

