📘 PANDUAN BIKIN TEMA BARU — Wedding Invitation Go
Versi: 3.0 (setelah Fase 2 + tema baru Botanical Garden)
Terakhir update: Setelah bikin tema Botanical Garden sebagai uji coba

🎯 Tujuan
Panduan ini untuk developer yang mau bikin tema undangan baru di project Wedding Invitation Go. Setelah baca ini, kamu bisa bikin tema baru yang:

Tampil di homepage (pilihan template)

Tampil di wizard Step 3 (pilih template)

Punya preview /preview/template/:slug

Berfungsi penuh (RSVP, music, footer, social links)

Cukup 1 file index.templ — tidak perlu helper, tidak perlu komponen manual

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
│   │   │   ├── web_footer.templ     ← web footer
│   │   │   └── social_links.templ   ← social links (icon + teks)
│   │   │
│   │   ├── rustic_wood/index.templ
│   │   ├── muslim_elegan/index.templ
│   │   ├── elegant_gold/index.templ
│   │   ├── modern_minimalist/index.templ
│   │   ├── traditional_java/index.templ
│   │   ├── botanical_garden/index.templ  ← tema baru (contoh)
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
     @shared.SocialLinks(data, shared.SocialLinksTextDefault)
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

❌ Social links (<div class="footer-social">...</div>)

GUNAKAN shared:

✅ @shared.RSVPForm(data, text) + @shared.RSVPListWrap()

✅ @shared.MusicPlayer()

✅ @shared.SocialLinks(data, text)

✅ @shared.WebFooter(data, text)

Kenapa? Supaya bug/update cukup di 1 tempat (bukan 6 tema).

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

⚠️ WAJIB panggil @shared.RSVPListWrap() SETELAH @shared.RSVPForm().

2. @shared.MusicPlayer()
Cara pakai (di bawah <body>, setelah overlay):

templ
@shared.MusicPlayer()
Parameter: tidak ada.

Prasyarat: <body> harus punya data-music-url, data-project-id, data-project-slug:

templ
<body data-music-url={ data.MusicURL } data-project-id={ invitation.FormatUint(data.Project.ID) } data-project-slug={ data.Project.Slug }>
3. @shared.SocialLinks(data, text)
Cara pakai (di dalam footer undangan):

templ
@shared.SocialLinks(data, shared.SocialLinksTextDefault)
Parameter text:

shared.SocialLinksTextDefault — "Instagram", "WhatsApp", "YouTube"

Otomatis menampilkan:

Instagram — icon + @handle (extract dari data.EventInstagram)

WhatsApp — icon + nomor (+62 812-3456-7890 dari data.Whatsapp)

YouTube — icon + URL (youtube.com dari data.Youtube)

Tombol hanya muncul kalau field tidak kosong. Kalau Whatsapp kosong → tombol WA hilang otomatis.

CSS wajib:

css
.footer-social {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
    margin-bottom: 24px;
}
.footer-social a {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 18px;
    border: 1px solid rgba(WARNA, 0.25);
    border-radius: 50px;
    color: rgb(WARNA);
    text-decoration: none;
    font-size: 0.8rem;
    font-weight: 500;
    transition: all 0.3s;
    opacity: 0.85;
    white-space: nowrap;
}
.footer-social a:hover {
    background: rgba(WARNA, 0.1);
    border-color: rgb(WARNA);
    transform: translateY(-2px);
    opacity: 1;
}
.footer-social a i { font-size: 1rem; }
.footer-social a span { font-size: 0.8rem; font-weight: 500; }
4. @shared.WebFooter(data, text)
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
.rsvp-form { max-width: 520px; margin: 0 auto; }
.form-group { margin-bottom: 20px; }
.form-group label { /* label */ }
.form-group input, .form-group select, .form-group textarea { /* input */ }
.form-group .required { color: #c0392b; }
.btn-submit { /* tombol submit */ }
.rsvp-warning { /* warning sudah RSVP */ }
.rsvp-warning-title { /* judul */ }
.rsvp-warning-title i { /* icon */ }
.rsvp-warning p { /* teks */ }
.rsvp-list-wrap { /* wrapper list */ }
.rsvp-list-head { /* head list */ }
.rsvp-list-head h4 { /* judul */ }
.rsvp-count { /* count badge */ }
.rsvp-empty { /* list kosong */ }
.rsvp-item { /* 1 item RSVP */ }
.rsvp-item .row-main { /* baris utama */ }
.rsvp-item .guest-name { /* nama tamu */ }
.rsvp-item .status-badge { /* badge status */ }
.rsvp-item .status-badge.hadir { background: rgba(16,185,129,0.1); color: #059669; }
.rsvp-item .status-badge.tidak_hadir { background: rgba(192,57,43,0.1); color: #c0392b; }
.rsvp-item .status-badge.ragu { background: rgba(245,158,11,0.12); color: #d97706; }
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
.web-footer-desc { /* deskripsi */ }
.web-footer-link { /* link */ }
.web-footer-socials { /* social container */ }
.web-social-link { /* social link */ }
.web-footer-heading { /* heading */ }
.web-footer-list { /* list */ }
.web-footer-version { /* versi */ }
.web-footer-divider { /* divider */ }
.web-footer-text { /* copyright */ }
⚠️ PENTING: 3 baris !important di atas WAJIB — kalau tidak:

.web-footer .container tanpa !important → footer tetap 1 kolom

.web-footer-brand img tanpa !important → logo & teks atas-bawah

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
JS standar yang bisa copy dari tema lain:
AOS init — AOS.init({...})

Music player — startMusic, toggleMusic, closeOverlay, initAudio, retry

Toast — showToast, escapeHtml

Copy bank — copyBankNumberFromEl, copyText, fallbackCopy

RSVP — loadRsvpList, submit handler

Lightbox — gallery click handler

Scroll helper — scrollToSection

Cara tercepat: Copy dari rustic_wood/index.templ, ganti nama class gallery & section ID.

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

            /* ============================================ */
            /* CSS TEMA (hero, couple, event, gallery, dll) */
            /* ============================================ */
            /* ... desain tema kamu di sini ... */

            /* ============================================ */
            /* RSVP FORM (WAJIB — copy dari rustic_wood) */
            /* ============================================ */
            .rsvp-form { max-width: 520px; margin: 0 auto; }
            .form-group { margin-bottom: 20px; }
            .form-group label { /* ... */ }
            .form-group input, .form-group select, .form-group textarea { /* ... */ }
            .form-group .required { color: #c0392b; }
            .btn-submit { /* ... */ }
            .rsvp-warning { /* ... */ }
            .rsvp-warning-title { /* ... */ }
            .rsvp-warning p { /* ... */ }
            .rsvp-list-wrap { /* ... */ }
            .rsvp-list-head { /* ... */ }
            .rsvp-count { /* ... */ }
            .rsvp-empty { /* ... */ }
            .rsvp-item { /* ... */ }
            /* dst. (lihat kontrak CSS di atas) */

            /* ============================================ */
            /* MUSIC PLAYER (WAJIB) */
            /* ============================================ */
            .music-player { position: fixed; bottom: 24px; right: 24px; z-index: 1000; }
            .music-btn { /* ... */ }
            .music-btn.playing { /* ... */ }

            /* ============================================ */
            /* FOOTER SOCIAL (WAJIB) */
            /* ============================================ */
            .footer-social { display: flex; justify-content: center; align-items: center; flex-wrap: wrap; gap: 10px; margin-bottom: 24px; }
            .footer-social a { display: inline-flex; align-items: center; gap: 8px; padding: 10px 18px; border: 1px solid rgba(WARNA, 0.25); border-radius: 50px; color: rgb(WARNA); text-decoration: none; font-size: 0.8rem; font-weight: 500; transition: all 0.3s; opacity: 0.85; white-space: nowrap; }
            .footer-social a:hover { background: rgba(WARNA, 0.1); border-color: rgb(WARNA); transform: translateY(-2px); opacity: 1; }
            .footer-social a i { font-size: 1rem; }
            .footer-social a span { font-size: 0.8rem; font-weight: 500; }

            /* ============================================ */
            /* WEB FOOTER (WAJIB) */
            /* ============================================ */
            .web-footer { /* ... */ }
            .web-footer .container { max-width: 1280px !important; }
            .web-footer-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; }
            .web-footer-brand img { display: inline-block !important; vertical-align: middle; margin-right: 8px; }
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
            <div class="overlay-inner">
                <!-- desain overlay -->
                <button class="btn-start" onclick={ templpkg.JSFuncCall("startMusic") }>
                    Buka Undangan
                </button>
            </div>
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

        <!-- FOOTER UNDANGAN (khusus tema) -->
        <footer class="footer-botanical">
            <!-- desain footer -->
            @shared.SocialLinks(data, shared.SocialLinksTextDefault)
        </footer>

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
⚠️ Cara tercepat: Copy dari rustic_wood/index.templ atau botanical_garden/index.templ, rename package & function, adjust CSS warna & ornamen.

Step 3: Register di render.go
Edit views/invitations/render.go:

go
import (
    "fmt"           // ← WAJIB
    "io"            // ← WAJIB
    "net/http"      // ← WAJIB
    // ⚠️ JANGAN import "context" — tidak dipakai, akan error

    "wedding-invitation-go/internal/invitation"
    "wedding-invitation-go/internal/models"
    "wedding-invitation-go/views/invitations/tema_baru"  // ← tambah
    // ... import tema lain
)

func RenderTemplate(...) error {
    // ...
    switch folder {
    case "rustic_wood": ...
    case "muslim_elegan": ...
    case "elegant_gold": ...
    case "modern_minimalist": ...
    case "traditional_java": ...
    case "botanical_garden": ...
    case "tema_baru":  // ← tambah case ini
        return tema_baru.TemaBaru(ctx, w, project, data, guestName)
    default:
        return fmt.Errorf("template %s tidak dikenal", folder)
    }
}
⚠️ WARNING: Jangan import "context" di render.go. RenderTemplate pakai r.Context(), tapi signature function tidak perlu import context (karena pakai *http.Request).

Step 4: Register di DB (template_seed.go)
Buka internal/database/template_seed.go, cari fungsi seedTemplates(), tambah di akhir sebelum return nil:

go
// ============================================
// 7. TEMA BARU
// ============================================
if err := upsertTemplate(
    "Tema Baru",               // Nama tampil
    "tema-baru",               // Slug (URL-friendly)
    "tema-baru",               // Folder (harus sama dengan nama folder)
    "Deskripsi singkat tema.", // Deskripsi
    buildSchema(defaultFields, library),
    7,                         // Order (7 = urutan ke-7)
); err != nil {
    return err
}
log.Println("✅ Template: Tema Baru")
Cara praktis pakai Python:

bash
cd /DATA/AppData/wedding-invitation-go/src

python3 << 'PYEOF'
path = "internal/database/template_seed.go"
with open(path, 'r') as f:
    content = f.read()

if 'TEMA BARU' in content:
    print("⚠️  Tema Baru sudah ada di seed")
else:
    new_block = """
	// ============================================
	// 7. TEMA BARU
	// ============================================
	if err := upsertTemplate(
		"Tema Baru",
		"tema-baru",
		"tema-baru",
		"Deskripsi singkat tema.",
		buildSchema(defaultFields, library),
		7,
	); err != nil {
		return err
	}
	log.Println("✅ Template: Tema Baru")

	return nil
}"""
    
    old = "\treturn nil\n}"
    idx = content.rfind(old)
    if idx != -1:
        content = content[:idx] + new_block + content[idx+len(old):]
        with open(path, 'w') as f:
            f.write(content)
        print("✅ Tema Baru ditambahkan ke seedTemplates()")
    else:
        print("❌ Pattern tidak ditemukan")
PYEOF
Step 5: Build & Restart
bash
cd /DATA/AppData/wedding-invitation-go/src

# Generate templ
docker exec -it wedding-invitation-go-app sh -c "cd /app && templ generate 2>&1 | tail -3"

# Build (WAJIB cek — kalau error, fix dulu)
docker exec -it wedding-invitation-go-app sh -c "cd /app && go build ./... 2>&1 | head -20"

# Restart (seed otomatis jalan)
docker restart wedding-invitation-go-app
sleep 15

# Cek log seed
docker logs wedding-invitation-go-app 2>&1 | grep -i "template" | tail -10
Step 6: Seed Template ke DB
Seed dijalankan otomatis saat startup (lihat main.go). Kalau tidak otomatis, cek DB:

bash
docker exec wedding-invitation-go-mysql mysql -uwedding_user -pwedding123 wedding_invitation_db -e "
SELECT id, name, slug, folder, is_active FROM templates ORDER BY id;
"
Yang diharapkan: ada row baru dengan slug = tema-baru, folder = tema-baru.

Kalau tidak ada, insert manual via admin panel (/admin/templates → Tambah Template) atau via DB.

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
File & Registrasi:

□ Folder views/invitations/tema_baru/ dibuat
□ File index.templ dibuat
□ func TemaBaru(ctx, w, project, data, guestName) error didefinisikan
□ templ temaBaruPage(data *invitation.TemplateData) didefinisikan
□ Import: shared, invitation, templpkg, models
□ <body> punya data-music-url, data-project-id, data-project-slug
Shared components:

□ @shared.MusicPlayer() dipanggil
□ @shared.RSVPForm(data, shared.RSVPTextIndonesia) dipanggil
□ @shared.RSVPListWrap() dipanggil SETELAH RSVPForm
□ @shared.SocialLinks(data, shared.SocialLinksTextDefault) dipanggil di footer undangan
□ @shared.WebFooter(data, shared.WebFooterTextDefault) dipanggil
CSS WAJIB:

□ .rsvp-form, .form-group, .btn-submit, .rsvp-warning, .rsvp-list-wrap
□ .music-player, .music-btn
□ .footer-social, .footer-social a, .footer-social a:hover, .footer-social a i, .footer-social a span
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

□ render.go diupdate (case baru, JANGAN import "context")
□ template_seed.go diupdate (upsertTemplate)
Build & Test:

□ templ generate sukses
□ go build ./... sukses
□ Seed template ke DB (auto via restart, atau manual)
□ Cek DB: SELECT * FROM templates WHERE slug='tema-baru';
□ Test preview /preview/template/tema-baru
□ Test wizard create project
□ Cek tampil di homepage
□ Cek footer undangan (social links)
□ Cek web footer (3 kolom)
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

Botanical: 🌿, 🍃, ✿, ❦

Test responsive — cek di mobile (max-width 600px) & tablet (600-1024px).

Copy dari tema lain — cara tercepat: copy rustic_wood/index.templ atau botanical_garden/index.templ, ganti nama package & function, adjust CSS warna & ornamen.

🚫 Yang TIDAK Boleh Diubah
❌ Jangan sentuh views/invitation/ — itu wizard, bukan template undangan

❌ Jangan edit file *_templ.go — itu generated, edit .templ saja

❌ Jangan pakai style={ background-image: ... } — templ escape ', pakai <img> tag

❌ Jangan pakai {{ var }} di dalam <script> — templ tidak replace {}, pakai data-* attribute

❌ Jangan tambah helpers.go per tema — semua helper global di internal/invitation/

❌ Jangan import "context" di render.go — tidak dipakai, akan error

❌ Jangan tulis form RSVP / music player / web footer / social links manual — pakai shared

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
Social links tidak muncul	Cek @shared.SocialLinks(data, shared.SocialLinksTextDefault) dipanggil + CSS .footer-social ada
Social links tumpang tindih	Cek .footer-social a { white-space: nowrap } + hapus CSS .footer-social a versi lama
CSS .footer-social dobel	Cari .footer-social a { — kalau ada 2, hapus yang versi lama (yang punya width: 40px)
color: 184, 149, 106 tidak valid	Ganti ke color: rgb(184, 149, 106)
Build error: "context" imported and not used	Hapus import "context" di render.go
Build error: "models" imported and not used	Hapus import "models" di shared component
Seed tidak jalan otomatis	Cek main.go — apakah panggil database.Seed()?
📞 Kontak & Referensi
Repo: https://github.com/falah-udin/wedding-invitation-go

Domain: https://wedding.litebox.my.id

Preview template: /preview/template/:slug

Contoh tema:

Paling sederhana: views/invitations/rustic_wood/index.templ

Contoh terbaru: views/invitations/botanical_garden/index.templ

📋 Ringkasan File yang Perlu Dibuat/Diubah
Untuk bikin tema baru:

views/invitations/tema_baru/index.templ — BARU (~500-1000 baris)

views/invitations/render.go — UPDATE (1 case + import)

internal/database/template_seed.go — UPDATE (1 upsertTemplate)

Total: 1 file baru + 2 file edit. 🎉

Estimasi waktu:

Copy dari tema lain: ~10 menit

Desain dari nol: ~2-4 jam

Testing: ~15 menit

Selamat berkarya! Kalau ada pertanyaan, tanya ke tim. 🚀

📌 Changelog
v3.0 (setelah bikin tema Botanical Garden):

Tambah shared/social_links.templ — komponen shared baru

Tambah section "SocialLinks" di Komponen Shared

Tambah CSS .footer-social di kontrak CSS

Tambah JS kontrak window.scrollToSection

Tambah warning context unused import di render.go

Tambah script Python praktis untuk register template_seed.go

Tambah debugging: .footer-social dobel, color: rgb() invalid, !important footer

Update checklist dengan SocialLinks

Tambah contoh tema terbaru botanical_garden

v2.0 (setelah Fase 2):

Tambah section "Komponen Shared — WAJIB DIPAKAI"

Tambah CSS wajib untuk web footer (.web-footer .container !important, .web-footer-brand img !important)

Tambah JS kontrak (startMusic, toggleMusic, scrollToSection)

Tambah 3 baris debugging footer (1 kolom, logo atas-bawah, dll)

Update "Yang TIDAK Boleh Diubah" — larangan tulis manual komponen shared

Update checklist dengan item shared components

v1.0 (setelah Fase 1):

Versi awal

