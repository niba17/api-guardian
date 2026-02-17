// web\dashboard-ui\src\features\dashboard\components\HistoryChart.tsx

import { useMemo } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import { CalendarClock } from "lucide-react"; // 👈 FIX: Import icon
import type { SecurityLog } from "../../../types/api.types";
import { formatTime } from "../../../utils/time";

// FIX: Definisikan interface untuk data chart
interface ChartDataPoint {
  time: string;
  valid: number;
  blocked: number;
}

export default function HistoryChart({ logs }: { logs: SecurityLog[] }) {
  const chartData = useMemo(() => {
    const grouped: Record<string, ChartDataPoint> = {};

    logs.forEach((log) => {
      // Kelompokkan per menit berdasarkan waktu lokal user
      const timeKey = formatTime(log.timestamp);

      if (!grouped[timeKey]) {
        grouped[timeKey] = { time: timeKey, valid: 0, blocked: 0 };
      }

      if (log.is_blocked) {
        grouped[timeKey].blocked++;
      } else {
        grouped[timeKey].valid++;
      }
    });

    return Object.values(grouped).sort((a, b) => a.time.localeCompare(b.time));
  }, [logs]);

  return (
    <div className="bg-slate-900/50 border border-slate-800 rounded-3xl p-6 h-88 flex flex-col">
      <div className="flex justify-between items-center mb-4 shrink-0">
        <div>
          <h3 className="font-bold text-lg text-white flex items-center gap-2">
            <CalendarClock className="text-orange-500" size={20} />
            Request History Volume
          </h3>
          <p className="text-slate-400 text-xs mt-1">
            Aggregated per minute (Session Persisted)
          </p>
        </div>
      </div>

      <div className="flex-1 w-full min-h-0">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={chartData}>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#1e293b"
              vertical={false}
            />
            <XAxis
              dataKey="time"
              stroke="#64748b"
              fontSize={11}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              stroke="#64748b"
              fontSize={11}
              tickLine={false}
              axisLine={false}
              allowDecimals={false}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#0f172a",
                borderColor: "#334155",
                borderRadius: "8px",
              }}
              itemStyle={{ fontSize: "12px" }}
              cursor={{ fill: "rgba(255,255,255,0.05)" }}
            />
            <Legend wrapperStyle={{ paddingTop: "10px", fontSize: "12px" }} />
            <Bar
              dataKey="valid"
              name="Valid Traffic"
              stackId="a"
              fill="#10b981"
              radius={[0, 0, 4, 4]}
              maxBarSize={50}
            />
            <Bar
              dataKey="blocked"
              name="Blocked Attacks"
              stackId="a"
              fill="#ef4444"
              radius={[4, 4, 0, 0]}
              maxBarSize={50}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
