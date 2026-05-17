#!/bin/sh
set -e

DB_FILE="${DB_PATH:-kare_rehber.db}"

# Sadece DB yoksa veya boşsa seed çalıştır
if [ ! -f "$DB_FILE" ] || [ "$(stat -c%s "$DB_FILE" 2>/dev/null || stat -f%z "$DB_FILE")" -lt 4096 ]; then
  echo "[start] DB bulunamadı, seed çalıştırılıyor..."
  ./seed
  echo "[start] Seed tamamlandı."
else
  echo "[start] Mevcut DB kullanılıyor: $DB_FILE"
fi

echo "[start] Sunucu başlatılıyor..."
exec ./server
