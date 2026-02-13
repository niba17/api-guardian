import { useState, useEffect } from "react";
import { api } from "../../../lib/axios";
import type { SystemStats, SecurityLog } from "../../../types/api.types";

export const useDashboardData = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [logs, setLogs] = useState<SecurityLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [lastUpdate, setLastUpdate] = useState<string>("");

  // --- LOGIKA PERSISTENCE ---

  // 1. Fungsi untuk mengambil history dari LocalStorage saat pertama kali load
  const getStoredLogs = (): SecurityLog[] => {
    try {
      const saved = localStorage.getItem("guardian_persistent_logs");
      return saved ? JSON.parse(saved) : [];
    } catch {
      // 👈 Cukup catch tanpa (e)
      return [];
    }
  };

  useEffect(() => {
    // Inisialisasi logs dengan data lama agar tidak kosong saat refresh
    const initialLogs = getStoredLogs();
    if (initialLogs.length > 0) {
      setLogs(initialLogs);
    }

    const fetchData = async () => {
      try {
        const resStats = await api.get("/dashboard/stats");
        setStats(resStats.data);

        const resLogs = await api.get("/dashboard/logs");
        const newIncomingLogs: SecurityLog[] = Array.isArray(resLogs.data)
          ? resLogs.data
          : [];

        // 2. MERGE DATA: Gabungkan data lama dan baru, buang duplikat berdasarkan ID
        setLogs((prevLogs) => {
          const logMap = new Map();

          // Masukkan data lama
          prevLogs.forEach((log) => logMap.set(log.id, log));

          // Masukkan data baru (akan menimpa jika ID sama, atau menambah jika baru)
          newIncomingLogs.forEach((log) => logMap.set(log.id, log));

          const merged = Array.from(logMap.values());

          // 3. AUTO-SAVE: Simpan hasil merge ke LocalStorage
          // Kita batasi 1000 logs saja supaya LocalStorage tidak bengkak
          const limitedLogs = merged.slice(-1000);
          localStorage.setItem(
            "guardian_persistent_logs",
            JSON.stringify(limitedLogs)
          );

          return limitedLogs;
        });

        setLastUpdate(new Date().toLocaleTimeString());
      } catch (error) {
        console.error("Backend Down / Error:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 2000);

    return () => clearInterval(interval);
  }, []);

  // Tambahkan fungsi manual clear jika dibutuhkan di UI
  const clearAllHistory = () => {
    localStorage.removeItem("guardian_persistent_logs");
    setLogs([]);
  };

  return { stats, logs, loading, lastUpdate, clearAllHistory };
};
