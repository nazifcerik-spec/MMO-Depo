# Animasyon öneri listesi

İlke (CLAUDE.md: kararlılık > performans > görsel): yalnız durum değişiminde bir kez oynayan, Flutter'ın yerleşik
araçlarıyla (AnimatedSwitcher, TweenAnimationBuilder, Hero, AnimatedContainer) yapılan, 150–600 ms'lik
animasyonlar. Sürekli döngü yok (dikkat gerektiren tek durum hariç), `MediaQuery.disableAnimations` açıksa hepsi
atlanır, yeni bağımlılık yok.

| # | Nerede | Tetik | Animasyon | Süre | Maliyet | Öncelik | Durum |
|---|---|---|---|---|---|---|---|
| 1 | Pet kartı | Pet evre atlar (Sv 34 / 67) | Eski evre görseli → yenisi çapraz geçiş + kısa ışık halkası | 600 ms | Düşük | Yüksek | Uygulandı — `SwapIn` + `GlowPulse` (pet kartı) |
| 2 | Pet kartı | Kuluçka bitti, yumurta hazır | Yumurta hafif sallanır (yalnız hazırken, 3 sn'de bir) | 400 ms | Düşük | Yüksek | Değişti — sürekli sallanma yerine kuluçka bitince yumurta → pet geçişi (`SwapIn`); döngü yok |
| 3 | Seviye atlama toast'ı | Karakter / meslek / pet seviye atlar | Rozet büyüyüp yerine oturur (scale 0.8→1.05→1) | 350 ms | Çok düşük | Yüksek | Uygulandı — ganimet sayfasında `PopIn` seviye rozeti; pet / görev toast'ı |
| 4 | XP / ilerleme çubukları | Değer değişir | Çubuk eski değerden yenisine dolar (zaten StatBar'da TweenAnimationBuilder ile) | 400 ms | Çok düşük | Yüksek | Uygulandı — `StatBar` (önceden animasyon yoktu) |
| 5 | Altın / XP sayıları (sonuç ekranı) | Ödül toplanır | Sayı 0'dan hedefe sayar | 500 ms | Çok düşük | Orta | Uygulandı — `CountUp` (ganimet sayfası altın + XP) |
| 6 | Eşya detayı | Çantada eşyaya dokunulur | Hero: hücredeki ikon detay başlığına uçar | 300 ms | Düşük | Orta | Ertelendi — eşya detayı alt sayfa (PopupRoute), Hero yalnız sayfa geçişinde uçar |
| 7 | Ganimet sayfası | Efsanevi / mitik eşya düşer | Nadirlik halesi bir kez parlar (mevcut aura, tek geçiş) | 600 ms | Düşük | Orta | Uygulandı — efsanevi / mitik satır `GlowPulse` |
| 8 | Efsunlama | Parşömen uygulanır | Eşya ikonunda sembol renginde tek parıltı | 500 ms | Düşük | Orta | Uygulandı — efsunlanan eşyanın hücresi (çanta / ekipman) `GlowPulse` |
| 9 | Zindan / baskın kartı | Toplamaya hazır | Madalyon kenarı bir kez nabız atar, sonra sabit vurgu | 800 ms | Düşük | Orta | Uygulandı — tüm süreli işler: `ContractProgress` çubuğu hazır olunca `GlowPulse` |
| 10 | Görev kartı | Hedef tamamlandı | Kart kenarı yeşile döner (AnimatedContainer) | 300 ms | Çok düşük | Orta | Uygulandı — görev kartı anahtarlandı, Material kenarı yerinde yeşile döner |
| 11 | Harita | Yolculuk başlar | Güzergâh çizgisi baştan sona çizilir | 800 ms | Düşük | Düşük | Uygulandı — ayrı `_RoutePainter` katmanı, rota baştan sona çizilir |
| 12 | Sekme geçişleri | Alt navigasyon | Yok (IndexedStack anında; ekleme önerilmez) | — | — | — | — |

Ortak araçlar `client/lib/shared/motion.dart` (DD-337): `motion()` azaltılmış harekette süreyi sıfırlar, `GlowPulse`, `SwapIn`, `PopIn`, `CountUp`; test `client/test/motion_test.dart`.

Önerilmeyenler: sürekli dönen/parlayan ikonlar, parçacık sistemleri, arka plan paralaksı, her liste satırına giriş
animasyonu — liste kaydırmasını ve pili yorar, akışı yavaşlatır.
