// web\dashboard-ui\src\utils\time.ts
import { formatInTimeZone } from "date-fns-tz"; // 👈 Pastikan sudah: npm install date-fns-tz

// 1. Ambil preferensi dari LocalStorage (jika user pernah ganti)
// 2. Jika tidak ada, pakai deteksi browser (default)
const DEFAULT_TIMEZONE =
  localStorage.getItem("user-timezone") ||
  Intl.DateTimeFormat().resolvedOptions().timeZone;

/**
 * Format string UTC dari database ke waktu lokal sesuai zona waktu terpilih
 */
export const formatTime = (utcString: string, pattern: string = "HH:mm:ss") => {
  try {
    // Pengaman: Jika string tidak punya penanda UTC, kita tambahkan 'Z'
    const cleanUtcString =
      utcString.endsWith("Z") || utcString.includes("+")
        ? utcString
        : `${utcString}Z`;

    const date = new Date(cleanUtcString);

    // Kita gunakan DEFAULT_TIMEZONE hasil pilihan user/auto
    return formatInTimeZone(date, DEFAULT_TIMEZONE, pattern);
  } catch (e) {
    console.error("Time conversion error:", e);
    return "Invalid Date";
  }
};

/**
 * Mendapatkan zona waktu yang sedang digunakan saat ini
 * Digunakan oleh TimezoneSelector untuk menampilkan value aktif
 */
export const getCurrentTimezone = () => DEFAULT_TIMEZONE;

/**
 * Mengatur zona waktu global, simpan ke storage, dan refresh halaman
 */
export const setGlobalTimezone = (tz: string) => {
  localStorage.setItem("user-timezone", tz);
  window.location.reload(); // Refresh total agar semua Chart & Table berubah serentak
};
