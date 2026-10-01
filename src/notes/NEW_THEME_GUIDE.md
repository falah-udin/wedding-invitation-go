📘 PANDUAN BIKIN TEMA BARU — Wedding Invitation Go

Versi: 3.1 (setelah Fase 2 + Botanical Garden + Dark Technology)
Terakhir update: Setelah bikin tema Dark Technology + fix bug overlay scroll
🎯 Tujuan

Panduan ini untuk developer yang mau bikin tema undangan baru di project Wedding Invitation Go. Setelah baca ini, kamu bisa bikin tema baru yang:

    Tampil di homepage (pilihan template)

    Tampil di wizard Step 3 (pilih template)

    Punya preview /preview/template/:slug

    Berfungsi penuh (RSVP, music, footer, social links, bank accounts)

    Anti-error — tidak akan kena bug templ parser

    Lock scroll overlay — user tidak bisa scroll di belakang overlay

🏗️ Arsitektur Project
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
│   │   └── preview_template.go      ← render preview template (data dummy)
│   │
│   ├── models/
│   │   ├── project.go
│   │   ├── template.go
│   │   └── bank_account.go
│   │
│   ├── services/
│   ├── database/
│   │   └── template_seed.go         ← SEED TEMPLATE + FIELD LIBRARY
│
├── views/
│   ├── invitations/                 ← TEMPLATE UNDANGAN
│   │   ├── shared/                  ← KOMPONEN SHARED
│   │   │   ├── rsvp.templ
│   │   │   ├── music_player.templ
│   │   │   ├── web_footer.templ
│   │   │   └── social_links.templ
│   │   │
│   │   ├── rustic_wood/
│   │   ├── muslim_elegan/
│   │   ├── elegant_gold/
│   │   ├── modern_minimalist/
│   │   ├── traditional_java/
│   │   ├── botanical_garden/
│   │   ├── dark_technology/
│   │   └── render.go                ← dispatcher (WAJIB DIUPDATE)
│   │
│   └── invitation/                  ← WIZARD (JANGAN DISENTUH)
│
└── public/storage/                  ← file upload & placeholder

🚨 ATURAN PALING PENTING (WAJIB BACA DULU)
1. ❌ JANGAN PAKAI // DI DALAM TAG HTML

Templ parser Go salah mengira // sebagai awal komentar, bikin error close tag not found.

❌ SALAH:
html

<div class="role">// The Groom</div>
<div class="hud">// BUILD v3.0.1</div>

✅ BENAR:
html

<div class="role">⟨ The Groom ⟩</div>
<div class="hud">[ BUILD v3.0.1 ]</div>
<div class="role"># The Groom</div>

Aturan praktis: Apapun teks di dalam <div>...</div>, jangan mulai dengan // atau ada // di posisi mana pun dalam tag.
2. ❌ JANGAN EDIT PAKAI sed BERULANG

sed -i 'NNNi\...' dan sed -i 'NNNd' berkali-kali = indentasi campur tab/spasi = templ parser gagal.

Penyebab utama cyber_dark gagal total selama 3 jam.

Gunakan:

    Python script dengan .replace() (multi-line safe)

    Editor manual (nano/vim) untuk edit kecil

    sed HANYA untuk rename string sederhana (1x pakai, bukan berulang)

3. ❌ HATI-HATI KARAKTER UNICODE EXOTIC

Beberapa karakter Unicode bikin parser Go bingung:
Aman ✅	Hindari ❌
◆ ◈ ● ◦	⎔ ◐ ◑ ◒ ◓
✦ ✿ ❦ ✧	⏚ ⏛ ⏜
⟨ ⟩ 「 」	⧉ ⧊ ⧋
↑ ↓ ← →	⇜ ⇝ ⇞

Kalau error "close tag not found" tapi tag seimbang, ganti karakter exotic.
4. ❌ JANGAN BUAT FILE DARI NOL

Cara paling aman bikin tema baru: copy dari tema yang sudah proven working.

    Copy dari botanical_garden atau dark_technology (paling baru & stabil)

    Rename package & function

    Baru styling bertahap

5. ❌ JANGAN LUPA RESTART CONTAINER

Setelah templ generate + go build, selalu docker restart wedding-invitation-go-app. Binary lama masih jalan kalau tidak restart.
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

Handler Preview Template (/preview/template/:slug)

Handler PreviewTemplate menggunakan data dummy hardcoded dari buildDummyData():

    show_bank_accounts: "yes" (sudah di-set)

    bank_accounts: [...] (4 rekening contoh)

    Foto, love stories, gallery — semua sudah di-set

Artinya: preview template selalu punya data lengkap. Kalau section bank tidak tampil di preview, masalahnya:

    File _templ.go belum di-regenerate → restart container

    Cache browser → Ctrl+Shift+R

    Bug di kode tema

📋 Struct TemplateData — Field yang Tersedia

Semua field diakses via data.FieldName.
Field Umum (dari data_undangan):
Field	Tipe	Keterangan
GroomName	string	Nama mempelai pria
BrideName	string	Nama mempelai wanita
GroomPhoto	string	Path foto pria
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
SiteConfig	models.Config	Setting situs (favicon, sosmed)
ExistingRsvp	*ExistingRsvp	RSVP tamu (kalau ada)
Specific	map[string]interface{}	Field dinamis dari template_specific_data
Project	models.Project	Project data (ID, Slug, dll)
GuestName	string	Nama tamu (dari ?to=)
IsNamedGuest	bool	true kalau tamu punya nama
MusicURL	string	URL musik (sudah full path)
Method Bantuan:
go

// Foto dengan fallback default
data.GroomPhotoOrDefault() string
data.BridePhotoOrDefault() string
data.FatherGroomPhotoOrDefault() string
data.MotherGroomPhotoOrDefault() string
data.FatherBridePhotoOrDefault() string
data.MotherBridePhotoOrDefault() string
data.HeroImageOr(fallback string) string
data.HasHeroImage() bool

// Resepsi label
data.ResepsiLabelDisplay() string
data.ResepsiArabic() string
data.IsIslamicResepsi() bool

// Orang tua
data.ParentsGroom() string
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
data.ExistingRsvpField(key, def string) string
data.ExistingRsvpGuests() int

Fungsi Statis (import "wedding-invitation-go/internal/invitation"):
go

// Bank
invitation.GetBankTypeBadge(acc) string
invitation.GetAccountBankName(acc) string
invitation.GetAccountInitial(acc) string
invitation.GetAccountNumber(acc) string
invitation.GetAccountName(acc) string
invitation.HasCustomIcon(acc) bool
invitation.GetAccountIconURL(acc) string

// Format
invitation.FormatUint(n uint) string
invitation.FormatInt(n int) string
invitation.CurrentYear() int
invitation.GuestGreeting(isNamed bool) string

📦 Komponen Shared — WAJIB DIPAKAI

⚠️ JANGAN tulis manual form RSVP, music player, web footer, social links. GUNAKAN shared.
1. @shared.RSVPForm(data, text) + @shared.RSVPListWrap()
templ

<section class="section-rsvp">
    <div class="container">
        <h2>Konfirmasi Kehadiran</h2>
        @shared.RSVPForm(data, shared.RSVPTextIndonesia)
        @shared.RSVPListWrap()
    </div>
</section>

⚠️ WAJIB panggil @shared.RSVPListWrap() SETELAH @shared.RSVPForm().
2. @shared.MusicPlayer()
templ

@shared.MusicPlayer()

Prasyarat: <body> harus punya data attribute:
templ

<body data-music-url={ data.MusicURL } data-project-id={ invitation.FormatUint(data.Project.ID) } data-project-slug={ data.Project.Slug }>

3. @shared.SocialLinks(data, text)
templ

@shared.SocialLinks(data, shared.SocialLinksTextDefault)

4. @shared.WebFooter(data, text)
templ

@shared.WebFooter(data, shared.WebFooterTextDefault)

📐 Struktur Section Standar (WAJIB IKUTI URUTAN)
text

1. Overlay (#startOverlay)         — khusus tema
2. @shared.MusicPlayer()            — WAJIB
3. Hero                             — khusus tema
4. Opening Quote                    — khusus tema
5. Couple                           — khusus tema
6. Love Story (if len > 0)          — opsional
7. Event                            — khusus tema
8. Gallery (if len > 0)             — opsional
9. RSVP (@shared.RSVPForm + List)   — WAJIB
10. Bank / Amplop Digital           — WAJIB kalau ada bank
11. Footer Undangan                 — khusus tema
12. @shared.WebFooter()             — WAJIB
13. Scripts (AOS, Swal, dll)

Section Bank (SERING TERLUPA!)
templ

<!-- AMPLOP DIGITAL -->
if data.ShouldShowBankAccounts() && len(data.BankAccounts) > 0 {
    <section class="section-{TEMA}-alt" id="section-bank">
        <div class="container">
            <div class="section-label" data-aos="fade-up">Amplop Digital</div>
            <h2 class="section-title" data-aos="fade-up" data-aos-delay="50">
                Kirim<br/><em>Hadiah</em>
            </h2>
            <p class="section-subtitle" data-aos="fade-up" data-aos-delay="100">
                Doa restu Anda adalah hadiah terindah. Namun jika ingin memberi, kami sediakan:
            </p>

            <div class="bank-{TEMA}-grid">
                for index, account := range data.BankAccounts {
                    <div class="bank-{TEMA}-card" data-aos="fade-up" data-aos-delay={ invitation.FormatInt(80 * (index + 1)) }>
                        <span class="bank-type-badge">{ invitation.GetBankTypeBadge(account) }</span>
                        <div class="bank-logo">
                            if invitation.HasCustomIcon(account) {
                                <img src={ invitation.GetAccountIconURL(account) } alt={ invitation.GetAccountBankName(account) }/>
                            } else {
                                <span class="bank-logo-icon">{ invitation.GetAccountInitial(account) }</span>
                            }
                        </div>
                        <div class="bank-name">{ invitation.GetAccountBankName(account) }</div>
                        <div class="bank-number-wrapper" data-number={ invitation.GetAccountNumber(account) } onclick={ templpkg.JSFuncCall("copyBankNumberFromEl", templpkg.JSExpression("this")) }>
                            <span class="bank-number">{ invitation.GetAccountNumber(account) }</span>
                            <i class="bi bi-clipboard copy-icon"></i>
                        </div>
                        <div class="bank-holder">
                            a.n. <strong>{ invitation.GetAccountName(account) }</strong>
                        </div>
                    </div>
                }
            </div>
        </div>
    </section>
}

🔒 Lock Body Scroll Saat Overlay Tampil (WAJIB)

Bug umum: overlay terbuka tapi body di belakang bisa di-scroll.
CSS (tambah di <style>):
css

body.overlay-active {
    overflow: hidden;
    position: fixed;
    width: 100%;
    height: 100%;
}

JavaScript (modifikasi closeOverlay()):
javascript

function closeOverlay() {
    var overlay = document.getElementById('startOverlay');
    if (overlay) {
        overlay.style.opacity = '0';
        setTimeout(function() {
            overlay.style.display = 'none';
        }, 800);
    }
    // Unlock scroll & paksa ke atas
    document.body.classList.remove('overlay-active');
    window.scrollTo(0, 0);
}

JavaScript init (di akhir <script>):
javascript

document.addEventListener('DOMContentLoaded', function() {
    var overlay = document.getElementById('startOverlay');
    if (overlay && overlay.style.display !== 'none') {
        document.body.classList.add('overlay-active');
    }
});

🎨 CSS — Kontrak Class
Class WAJIB (shared mengandalkan ini):
css

/* RSVP */
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
.rsvp-item .row-main { /* ... */ }
.rsvp-item .guest-name { /* ... */ }
.rsvp-item .status-badge { /* ... */ }
.rsvp-item .status-badge.hadir { background: rgba(16,185,129,0.1); color: #059669; }
.rsvp-item .status-badge.tidak_hadir { background: rgba(192,57,43,0.1); color: #c0392b; }
.rsvp-item .status-badge.ragu { background: rgba(245,158,11,0.12); color: #d97706; }
.rsvp-item .row-detail { /* ... */ }
.rsvp-item .meta-left { /* ... */ }
.rsvp-item .message { /* ... */ }

/* Music Player */
.music-player { position: fixed; bottom: 24px; right: 24px; z-index: 1000; }
.music-btn { /* ... */ }
.music-btn.playing { /* ... */ }

/* Footer Social */
.footer-social { display: flex; justify-content: center; align-items: center; flex-wrap: wrap; gap: 10px; margin-bottom: 24px; }
.footer-social a { display: inline-flex; align-items: center; gap: 8px; padding: 10px 18px; border: 1px solid rgba(WARNA, 0.25); border-radius: 50px; color: rgb(WARNA); text-decoration: none; font-size: 0.8rem; font-weight: 500; transition: all 0.3s; opacity: 0.85; white-space: nowrap; }
.footer-social a:hover { background: rgba(WARNA, 0.1); border-color: rgb(WARNA); transform: translateY(-2px); opacity: 1; }
.footer-social a i { font-size: 1rem; }
.footer-social a span { font-size: 0.8rem; font-weight: 500; }

/* Web Footer (3 baris !important WAJIB) */
.web-footer .container { max-width: 1280px !important; }
.web-footer-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; }
.web-footer-brand img { display: inline-block !important; vertical-align: middle; margin-right: 8px; }
.web-footer-icon { color: var(--warna-tema); }

/* Toast */
.toast-container { position: fixed; top: 20px; right: 20px; z-index: 9999; }
.toast { /* ... */ }

/* Overlay Scroll Lock */
body.overlay-active { overflow: hidden; position: fixed; width: 100%; height: 100%; }

🛠️ CARA BIKIN TEMA BARU (Step-by-Step)
Step 1: Copy dari Tema Proven
bash

cd /DATA/AppData/wedding-invitation-go/src/views/invitations

# Copy dari botanical_garden atau dark_technology
cp -r botanical_garden tema_baru

# Hapus file generated
rm -f tema_baru/index_templ.go
rm -f tema_baru/*.bak*

Step 2: Rename Package & Function

Pakai Python (aman multi-line):
bash

python3 << 'PYEOF'
path = "/DATA/AppData/wedding-invitation-go/src/views/invitations/tema_baru/index.templ"
with open(path) as f:
    content = f.read()

content = content.replace('package botanical_garden', 'package tema_baru')
content = content.replace('func BotanicalGarden(', 'func TemaBaru(')
content = content.replace('botanicalGardenPage(', 'temaBaruPage(')
content = content.replace('- Botanical Garden</title>', '- Tema Baru</title>')

with open(path, 'w') as f:
    f.write(content)
print("✅ Rename selesai")
PYEOF

Step 3: Verifikasi
bash

docker exec -it wedding-invitation-go-app sh -c "head -1 /app/views/invitations/tema_baru/index.templ"
docker exec -it wedding-invitation-go-app sh -c "grep -n 'func TemaBaru\|temaBaruPage' /app/views/invitations/tema_baru/index.templ"

Step 4: Generate Awal (WAJIB sukses sebelum lanjut styling)
bash

docker exec -it wedding-invitation-go-app sh -c "cd /app && templ generate 2>&1 | grep -iE 'tema_baru|error' | head -5"
docker exec -it wedding-invitation-go-app sh -c "cd /app && go build ./... 2>&1 | head -10"

⚠️ Kalau generate error di step ini, JANGAN LANJUT styling. Fix dulu.
Step 5: Register di render.go
go

import (
    // ... import lain
    "wedding-invitation-go/views/invitations/tema_baru"
)

func RenderTemplate(...) error {
    switch folder {
    // ... case lain
    case "tema_baru":
        return tema_baru.TemaBaru(ctx, w, project, data, guestName)
    default:
        return fmt.Errorf("template %s tidak dikenal", folder)
    }
}

⚠️ JANGAN import "context" di render.go — tidak dipakai, akan error.
Step 6: Register di template_seed.go

Cari seedTemplates() di /app/internal/database/template_seed.go, tambah sebelum return nil:
go

// ============================================
// N. TEMA BARU
// ============================================
if err := upsertTemplate(
    "Tema Baru",
    "tema-baru",
    "tema-baru",
    "Deskripsi singkat tema.",
    buildSchema(defaultFields, library),
    N,   // order
); err != nil {
    return err
}
log.Println("✅ Template: Tema Baru")

Cara aman pakai Python:
bash

python3 << 'PYEOF'
path = "/DATA/AppData/wedding-invitation-go/src/internal/database/template_seed.go"
with open(path) as f:
    content = f.read()

if 'Template: Tema Baru' in content:
    print("⚠️  Tema Baru sudah ada di seed")
else:
    new_block = '''        // ============================================
        // N. TEMA BARU
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
}'''
    
    old = "        return nil\n}"
    idx = content.rfind(old)
    if idx == -1:
        print("❌ Pattern tidak ditemukan")
    else:
        content = content[:idx] + new_block + content[idx+len(old):]
        with open(path, 'w') as f:
            f.write(content)
        print("✅ Tema Baru ditambahkan ke seedTemplates()")
PYEOF

Step 7: Build & Restart
bash

docker exec -it wedding-invitation-go-app sh -c "cd /app && templ generate 2>&1 | tail -3"
docker exec -it wedding-invitation-go-app sh -c "cd /app && go build ./... 2>&1 | head -20"

# WAJIB restart supaya seed jalan
docker restart wedding-invitation-go-app
sleep 15

# Cek log seed
docker logs wedding-invitation-go-app 2>&1 | grep -i "template" | tail -10

Step 8: Verifikasi di DB
bash

docker exec wedding-invitation-go-mysql mysql -uwedding_user -pwedding123 wedding_invitation_db -e "
SELECT id, name, slug, folder, \`order\`, is_active FROM templates ORDER BY \`order\`;
"

Harus ada row baru dengan slug = tema-baru.
Step 9: Test Preview

Buka: https://wedding.litebox.my.id/preview/template/tema-baru
Step 10: Styling Bertahap

Patch satu-satu:

    Ganti CSS variables (:root) → generate → build → restart → cek

    Ganti font link → generate → build → restart → cek

    Ganti ornamen (emoji) → generate → build → restart → cek

    Ganti class name (.botanical-* → .tema-*) → generate → build → restart → cek

    Ganti teks & judul section → generate → build → restart → cek

Setiap patch: generate → build → restart → refresh browser.

⚠️ JANGAN pakai sed berulang. Pakai Python .replace().
✅ Checklist Tema Baru
File & Registrasi:

    □

    Folder views/invitations/tema_baru/ dibuat
    □

    File index.templ dibuat (copy dari tema proven)
    □

    func TemaBaru(ctx, w, project, data, guestName) error didefinisikan
    □

    templ temaBaruPage(data *invitation.TemplateData) didefinisikan
    □

    Import: shared, invitation, templpkg, models
    □

    <body> punya data-music-url, data-project-id, data-project-slug

Anti-Pattern (WAJIB CEK):

    □

    Tidak ada // di dalam tag HTML
    □

    Tidak ada karakter Unicode exotic (⎔ ◐ ◑ ◒ ◓)
    □

    Indentasi konsisten (tab atau spasi, jangan campur)
    □

    Tidak ada nested <div> dengan indentasi mundur di dalam <if>

Shared Components:

    □

    @shared.MusicPlayer() dipanggil
    □

    @shared.RSVPForm(data, shared.RSVPTextIndonesia) dipanggil
    □

    @shared.RSVPListWrap() dipanggil SETELAH RSVPForm
    □

    @shared.SocialLinks(data, shared.SocialLinksTextDefault) dipanggil di footer
    □

    @shared.WebFooter(data, shared.WebFooterTextDefault) dipanggil

Section Bank (WAJIB):

    □

    Section "Amplop Digital" ada (setelah RSVP, sebelum Footer)
    □

    Pakai if data.ShouldShowBankAccounts() && len(data.BankAccounts) > 0
    □

    Pakai invitation.GetBankTypeBadge() dll

Lock Scroll Overlay:

    □

    CSS body.overlay-active { overflow: hidden; position: fixed; }
    □

    JS closeOverlay() add classList.remove('overlay-active') + window.scrollTo(0, 0)
    □

    JS init saat DOMContentLoaded add overlay-active class

CSS WAJIB:

    □

    .rsvp-form, .form-group, .btn-submit, .rsvp-warning, .rsvp-list-wrap
    □

    .music-player, .music-btn
    □

    .footer-social, .footer-social a, :hover, a i, a span
    □

    .web-footer, .web-footer-row
    □

    .web-footer .container { max-width: 1280px !important }
    □

    .web-footer-brand img { display: inline-block !important; vertical-align: middle }
    □

    .web-footer-icon { color: var(--warna-tema) }
    □

    .toast-container, .toast
    □

    body.overlay-active

JS WAJIB:

    □

    window.startMusic, window.toggleMusic
    □

    window.scrollToSection (atau alias)
    □

    showToast, escapeHtml
    □

    copyBankNumberFromEl, copyText, fallbackCopy
    □

    loadRsvpList + submit handler

Registrasi:

    □

    render.go diupdate (case baru, JANGAN import "context")
    □

    template_seed.go diupdate (upsertTemplate)

Build & Test:

    □

    templ generate sukses
    □

    go build ./... sukses
    □

    Container direstart
    □

    Seed template ke DB (auto via restart)
    □

    Cek DB: SELECT * FROM templates WHERE slug='tema-baru';
    □

    Test preview /preview/template/tema-baru
    □

    Test wizard create project
    □

    Cek tampil di homepage
    □

    Cek footer undangan (social links)
    □

    Cek web footer (3 kolom)
    □

    Cek section bank muncul (Amplop Digital)
    □

    Cek overlay: body tidak bisa scroll saat overlay terbuka
    □

    Cek overlay: klik → mulai dari hero (bukan di tengah)

🐛 Debugging — Problem & Solusi
Problem	Solusi
templ generate error "close tag not found"	Cek // di dalam tag HTML — paling sering!
Error "close tag not found" tapi tag seimbang	Ganti karakter Unicode exotic (⎔ ◐ → ◆ ●)
Error "close tag not found" berulang	Indentasi campur tab/spasi — pakai Python replace
templ generate skip file	Hapus _templ.go, generate ulang
Field tidak tampil	Cek data.FieldName ada di struct? Cek LoadData
Preview 500 error	docker logs wedding-invitation-go-app
Template tidak muncul di homepage	SELECT * FROM templates WHERE is_active=1;
Preview blank	Cek error JS di browser console (F12)
Gambar tidak muncul	Cek path /storage/... atau URL eksternal
Music tidak play	Cek data-music-url di <body>, cek console
Footer 1 kolom, bukan 3	Cek .web-footer .container { max-width: 1280px !important }
Logo + teks footer atas-bawah	Cek .web-footer-brand img { display: inline-block !important }
RSVP form hilang	Cek @shared.RSVPForm + @shared.RSVPListWrap dipanggil
Music player tidak muncul	Cek @shared.MusicPlayer() dipanggil
Social links tidak muncul	Cek @shared.SocialLinks(...) dipanggil + CSS .footer-social ada
Section bank tidak muncul di preview	Restart container (file _templ.go lama)
Section bank tidak muncul di undangan	Cek data project: show_bank_accounts = "yes"
Body bisa di-scroll saat overlay terbuka	Tambah body.overlay-active CSS + JS closeOverlay
Setelah klik overlay, mulai di tengah	Tambah window.scrollTo(0, 0) di closeOverlay
sed bikin file rusak	Pakai Python .replace() multi-line
Build error: "context" imported and not used	Hapus import "context" di render.go
Seed tidak jalan otomatis	Cek main.go panggil database.Seed()
🚫 Yang TIDAK Boleh Diubah

    ❌ Jangan sentuh views/invitation/ — itu wizard

    ❌ Jangan edit file *_templ.go — itu generated

    ❌ Jangan pakai style={ background-image: ... } — pakai <img> tag

    ❌ Jangan pakai {{ var }} di dalam <script> — pakai data-* attribute

    ❌ Jangan tambah helpers.go per tema

    ❌ Jangan import "context" di render.go

    ❌ Jangan tulis manual form RSVP / music / footer / social links — pakai shared

    ❌ Jangan lupa @shared.RSVPListWrap() setelah @shared.RSVPForm()

    ❌ Jangan pakai // di dalam tag HTML ← BARU

    ❌ Jangan pakai sed berulang untuk edit file .templ ← BARU

    ❌ Jangan pakai karakter Unicode exotic ⎔ ◐ ◑ ◒ ◓ ← BARU

    ❌ Jangan buat file dari nol — copy dari tema proven ← BARU

📞 Kontak & Referensi

    Repo: https://github.com/falah-udin/wedding-invitation-go

    Domain: https://wedding.litebox.my.id

    Preview template: /preview/template/:slug

Contoh tema:

    Paling sederhana: views/invitations/rustic_wood/index.templ

    Modern: views/invitations/modern_minimalist/index.templ

    Terbaru (dengan HUD & terminal): views/invitations/dark_technology/index.templ

    Dengan section bank lengkap: views/invitations/elegant_gold/index.templ

📋 Ringkasan File yang Perlu Dibuat/Diubah

Untuk bikin tema baru:

    views/invitations/tema_baru/index.templ — BARU

    views/invitations/render.go — UPDATE (1 case + import)

    internal/database/template_seed.go — UPDATE (1 upsertTemplate)

Total: 1 file baru + 2 file edit.
Estimasi waktu:

    Copy dari tema proven: ~10 menit

    Styling dari nol: ~2-4 jam

    Testing: ~15 menit

📌 Changelog
v3.1 (setelah bikin Dark Technology)

Insight baru dari sesi dark_technology:

    🔴 Tambah aturan KERAS: Jangan pakai // di dalam tag HTML → penyebab error "close tag not found"

    🔴 Tambah aturan: Jangan edit pakai sed berulang — pakai Python .replace()

    🔴 Tambah aturan: Hindari karakter Unicode exotic (⎔ ◐ ◑ ◒ ◓)

    🔴 Tambah aturan: Copy dari tema proven, jangan bikin dari nol

    🔴 Tambah aturan: Selalu restart container setelah templ generate + go build

    🔴 Tambah section: "Lock Body Scroll Saat Overlay Tampil" (fix bug umum)

    🔴 Tambah section: "Struktur Section Standar" dengan urutan 1-13

    🔴 Tambah section: "Section Bank / Amplop Digital" (sering terlupa)

    🔴 Tambah catatan: Handler PreviewTemplate pakai data dummy — kalau bank tidak muncul, restart container

    🟡 Update Debugging: 6 baris baru (termasuk "section bank tidak muncul di preview")

    🟡 Update Checklist: 5 item baru (anti-pattern, section bank, lock overlay)

    🟡 Update Yang TIDAK Boleh Diubah: 4 item baru

v3.0 (setelah bikin Botanical Garden)

    Tambah shared/social_links.templ

    Tambah section "SocialLinks" di Komponen Shared

    Tambah CSS .footer-social di kontrak CSS

    Tambah JS kontrak window.scrollToSection

    Tambah warning context unused import di render.go

    Tambah debugging: .footer-social dobel, color: rgb() invalid, !important footer

v2.0 (setelah Fase 2)

    Tambah section "Komponen Shared — WAJIB DIPAKAI"

    Tambah CSS wajib untuk web footer

    Tambah JS kontrak

v1.0 (setelah Fase 1)

    Versi awal

📌 Catatan Tambahan — Favicon Dinamis

Setiap tema undangan wajib punya favicon dinamis. Sisipkan di index.templ setelah </title>:
templ

<title>{ data.GroomName } &amp; { data.BrideName } - Nama Tema</title>

<!-- Favicon Dinamis -->
if data.SiteConfig.SiteFavicon != "" {
    <link rel="icon" href={ "/storage/" + data.SiteConfig.SiteFavicon }/>
    <link rel="apple-touch-icon" href={ "/storage/" + data.SiteConfig.SiteFavicon }/>
} else {
    <link rel="icon" type="image/x-icon" href="/favicon.ico"/>
}

Cek semua tema punya favicon:
bash

cd /DATA/AppData/wedding-invitation-go/src
for tema in rustic_wood muslim_elegan elegant_gold modern_minimalist traditional_java botanical_garden dark_technology; do
  echo "--- $tema ---"
  grep -c "SiteConfig.SiteFavicon" views/invitations/$tema/index.templ
done

Output harusnya 1 di setiap tema.
🚀 Selamat Berkarya!

Kalau ada pertanyaan, tanya ke tim. Ingat aturan emas:

    Copy dari tema proven — jangan bikin dari nol

    Jangan pakai // di HTML

    Generate + build + restart setiap patch

    Test preview setelah setiap perubahan

Happy coding! 🎉