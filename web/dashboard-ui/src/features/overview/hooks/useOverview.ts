import { useState, useEffect } from "react";
import api from "../../../lib/axios";
import type { SystemStats, SecurityLog } from "../../../types/api.types";

export const useOverviewData = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [logs, setLogs] = useState<SecurityLog[]>([]); // Default array kosong
  const [loading, setLoading] = useState(true);
  const [lastUpdate, setLastUpdate] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [resStats, resLogs] = await Promise.all([
          api.get("/dashboard/stats"),
          api.get("/dashboard/logs"),
        ]);

        setStats(resStats.data);

        // Langsung set data dari server, tidak perlu merge/localstorage
        // Server sudah menjamin urutan dan kelengkapan data
        setLogs(Array.isArray(resLogs.data) ? resLogs.data : []);

        setLastUpdate(new Date().toLocaleTimeString());
      } catch (error) {
        console.error("Overview Fetch Error:", error);
      } finally {
        setLoading(false);
      }
    };

    // 1. Panggil pertama kali saat load
    fetchData();

    // 2. Polling tiap 2 detik untuk update real-time
    const interval = setInterval(fetchData, 2000);

    return () => clearInterval(interval);
  }, []);

  // Tidak perlu lagi function 'clearAllHistory' karena data ada di server

  return { stats, logs, loading, lastUpdate };
};
