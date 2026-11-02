# Animasyon öneri listesi

İlke (CLAUDE.md: kararlılık > performans > görsel): yalnız durum değişiminde bir kez oynayan, Flutter'ın yerleşik
araçlarıyla (AnimatedSwitcher, TweenAnimationBuilder, Hero, AnimatedContainer) yapılan, 150–600 ms'lik
animasyonlar. Sürekli döngü yok (dikkat gerektiren tek durum hariç), `MediaQuery.disableAnimations` açıksa hepsi
atlanır, yeni bağımlılık yok.

| # | Nerede | Tetik | Animasyon | Süre | Maliyet | Öncelik |
|---|---|---|---|---|---|---|
| 1 | Pet kartı | Pet evre atlar (Sv 34 / 67) | Eski evre görseli → yenisi çapraz geçiş + kısa ışık halkası | 600 ms | Düşük | Yüksek |
| 2 | Pet kartı | Kuluçka bitti, yumurta hazır | Yumurta hafif sallanır (yalnız hazırken, 3 sn'de bir) | 400 ms | Düşük | Yüksek |
| 3 | Seviye atlama toast'ı | Karakter / meslek / pet seviye atlar | Rozet büyüyüp yerine oturur (scale 0.8→1.05→1) | 350 ms | Çok düşük | Yüksek |
| 4 | XP / ilerleme çubukları | Değer değişir | Çubuk eski değerden yenisine dolar (zaten StatBar'da TweenAnimationBuilder ile) | 400 ms | Çok düşük | Yüksek |
| 5 | Altın / XP sayıları (sonuç ekranı) | Ödül toplanır | Sayı 0'dan hedefe sayar | 500 ms | Çok düşük | Orta |
| 6 | Eşya detayı | Çantada eşyaya dokunulur | Hero: hücredeki ikon detay başlığına uçar | 300 ms | Düşük | Orta |
| 7 | Ganimet sayfası | Efsanevi / mitik eşya düşer | Nadirlik halesi bir kez parlar (mevcut aura, tek geçiş) | 600 ms | Düşük | Orta |
| 8 | Efsunlama | Parşömen uygulanır | Eşya ikonunda sembol renginde tek parıltı | 500 ms | Düşük | Orta |
| 9 | Zindan / baskın kartı | Toplamaya hazır | Madalyon kenarı bir kez nabız atar, sonra sabit vurgu | 800 ms | Düşük | Orta |
| 10 | Görev kartı | Hedef tamamlandı | Kart kenarı yeşile döner (AnimatedContainer) | 300 ms | Çok düşük | Orta |
| 11 | Harita | Yolculuk başlar | Güzergâh çizgisi baştan sona çizilir | 800 ms | Düşük | Düşük |
| 12 | Sekme geçişleri | Alt navigasyon | Yok (IndexedStack anında; ekleme önerilmez) | — | — | — |

Önerilmeyenler: sürekli dönen/parlayan ikonlar, parçacık sistemleri, arka plan paralaksı, her liste satırına giriş
animasyonu — liste kaydırmasını ve pili yorar, akışı yavaşlatır.
