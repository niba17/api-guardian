import { useState, useEffect, useMemo } from "react";
import api from "../../../lib/axios";
import type { SecurityLog } from "../../../types/api.types";

export const useWAFData = () => {
  const [logs, setLogs] = useState<SecurityLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const resLogs = await api.get("/dashboard/logs");

        // 🚀 Bersih dan langsung dari server
        setLogs(Array.isArray(resLogs.data) ? resLogs.data : []);
      } catch (error) {
        console.error("WAF Data Fetch Error:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 2000);
    return () => clearInterval(interval);
  }, []);

  // --- DERIVED DATA (Logika Bisnis WAF) ---
  const wafStats = useMemo(() => {
    // 1. Filter hanya log yang DIBLOKIR
    const threats = logs.filter((log) => log.is_blocked);

    // 2. Hitung Unique Countries
    const uniqueCountries = new Set(threats.map((l) => l.country)).size;

    // 3. Cari Target Terpopuler (Endpoint yang paling sering diserang)
    const targetCounts = threats.reduce((acc, log) => {
      acc[log.path] = (acc[log.path] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);

    const topTarget =
      Object.entries(targetCounts).sort((a, b) => b[1] - a[1])[0]?.[0] || "-";

    return {
      threats, // Array log serangan
      totalThreats: threats.length,
      uniqueCountries,
      topTarget,
    };
  }, [logs]);

  return {
    loading,
    ...wafStats, // Spread stats biar langsung bisa dipakai (threats, totalThreats, etc)
  };
};
