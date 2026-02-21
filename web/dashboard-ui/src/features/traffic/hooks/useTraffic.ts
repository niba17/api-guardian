import { useState, useEffect, useMemo } from "react";
import api from "../../../lib/axios";
import type { SystemStats, SecurityLog } from "../../../types/api.types";

export const useTrafficData = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [logs, setLogs] = useState<SecurityLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [resStats, resLogs] = await Promise.all([
          api.get("/dashboard/stats"),
          api.get("/dashboard/logs"),
        ]);

        setStats(resStats.data);

        // 🚀 Langsung set data dari server, bersih tanpa merge localStorage!
        setLogs(Array.isArray(resLogs.data) ? resLogs.data : []);
      } catch (error) {
        console.error("Traffic Data Fetch Error:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 2000);
    return () => clearInterval(interval);
  }, []);

  // --- DERIVED LOGIC: Top Endpoints ---
  const topEndpoints = useMemo(() => {
    const counts = logs.reduce((acc, log) => {
      acc[log.path] = (acc[log.path] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);

    return Object.entries(counts)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5);
  }, [logs]);

  return { stats, logs, loading, topEndpoints };
};
