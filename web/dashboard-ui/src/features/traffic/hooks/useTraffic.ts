import { useState, useEffect, useMemo } from "react";
import api from "../../../lib/axios";
import type { SystemStats, SecurityLog } from "../../../types/api.types";

export const useTrafficData = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [logs, setLogs] = useState<SecurityLog[]>([]);
  const [loading, setLoading] = useState(true);

  // --- PERSISTENCE LOGIC ---
  const getStoredLogs = (): SecurityLog[] => {
    try {
      const saved = localStorage.getItem("guardian_persistent_logs");
      return saved ? JSON.parse(saved) : [];
    } catch {
      return [];
    }
  };

  useEffect(() => {
    const initialLogs = getStoredLogs();
    if (initialLogs.length > 0) setLogs(initialLogs);

    const fetchData = async () => {
      try {
        const [resStats, resLogs] = await Promise.all([
          api.get("/dashboard/stats"),
          api.get("/dashboard/logs"),
        ]);

        setStats(resStats.data);
        const newIncomingLogs: SecurityLog[] = Array.isArray(resLogs.data)
          ? resLogs.data
          : [];

        setLogs((prevLogs) => {
          const logMap = new Map();
          prevLogs.forEach((log) => logMap.set(log.id, log));
          newIncomingLogs.forEach((log) => logMap.set(log.id, log));

          const merged = Array.from(logMap.values());
          const limitedLogs = merged.slice(-1000);

          localStorage.setItem(
            "guardian_persistent_logs",
            JSON.stringify(limitedLogs)
          );
          return limitedLogs;
        });
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
