import { useMemo } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { SecurityLog } from "../../../types/api.types";

interface LatencyChartProps {
  logs: SecurityLog[];
}

export default function LatencyChart({ logs }: LatencyChartProps) {
  const data = useMemo(() => {
    // Ambil 20 log terakhir biar grafik gak terlalu padat
    return logs
      .slice(0, 20)
      .reverse()
      .map((log) => ({
        time: new Date(log.timestamp).toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        }),
        latency: log.latency,
        status: log.status,
      }));
  }, [logs]);

  return (
    <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6 h-80">
      <div className="flex justify-between items-center mb-4">
        <h3 className="font-bold text-lg text-white">System Latency (ms)</h3>
        <span className="text-xs px-2 py-1 bg-slate-800 rounded text-slate-400">
          Live (Last 20 Reqs)
        </span>
      </div>

      <ResponsiveContainer width="100%" height="80%">
        <LineChart data={data}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="#1e293b"
            vertical={false}
          />
          <XAxis
            dataKey="time"
            stroke="#64748b"
            fontSize={10}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            stroke="#64748b"
            fontSize={10}
            tickLine={false}
            axisLine={false}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: "#0f172a",
              borderColor: "#334155",
              color: "#fff",
            }}
          />
          {/* Garis Kuning Petir */}
          <Line
            type="monotone"
            dataKey="latency"
            stroke="#facc15"
            strokeWidth={3}
            dot={{ r: 4, fill: "#facc15" }}
            activeDot={{ r: 6 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
