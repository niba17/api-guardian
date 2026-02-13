import { useState, useEffect, useRef } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Activity } from "lucide-react";
import type { SecurityLog } from "../../../types/api.types";

interface TrafficChartProps {
  logs: SecurityLog[];
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{ value: number; stroke: string }>;
  label?: string;
}

const CustomTooltip = ({ active, payload, label }: CustomTooltipProps) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-slate-900 border border-slate-700 p-3 rounded-lg shadow-xl backdrop-blur-md">
        <p className="text-slate-400 text-xs font-mono mb-2 border-b border-slate-700 pb-1">
          {label}
        </p>
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
            <p className="text-sm font-bold text-slate-200">
              Valid:{" "}
              <span className="text-emerald-400">{payload[0]?.value}</span>
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-red-500"></span>
            <p className="text-sm font-bold text-slate-200">
              Blocked: <span className="text-red-400">{payload[1]?.value}</span>
            </p>
          </div>
        </div>
      </div>
    );
  }
  return null;
};

export default function TrafficChart({ logs }: TrafficChartProps) {
  // 1. STATE GRAFIK
  const [dataHistory, setDataHistory] = useState<
    Array<{ time: string; valid: number; blocked: number }>
  >([]);

  // 2. REF UNTUK MENYIMPAN LOG TERBARU (Tanpa Memicu Render)
  // Ini trik ninja supaya useEffect timer bisa baca data terbaru tanpa restart timer
  const latestLogsRef = useRef<SecurityLog[]>([]);

  // Update ref setiap kali props logs berubah
  useEffect(() => {
    latestLogsRef.current = logs;
  }, [logs]);

  // 3. TIMER PENGGERAK GRAFIK (Jantung Detak Grafik)
  useEffect(() => {
    const updateChart = () => {
      const now = new Date();
      const timeLabel = now.toLocaleTimeString("en-GB", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });

      // Ambil data dari REF, bukan dari state/props langsung
      // Filter log yang umurnya < 5 detik dari sekarang
      const recentLogs = latestLogsRef.current.filter((log) => {
        const logTime = new Date(log.timestamp).getTime();
        return now.getTime() - logTime < 5000;
      });

      const validCount = recentLogs.filter((l) => !l.is_blocked).length;
      const blockedCount = recentLogs.filter((l) => l.is_blocked).length;

      setDataHistory((prev) => {
        const newData = [
          ...prev,
          { time: timeLabel, valid: validCount, blocked: blockedCount },
        ];
        // Jaga agar tetap 20 titik data
        if (newData.length > 20) {
          return newData.slice(newData.length - 20);
        }
        return newData;
      });
    };

    // Jalankan updateChart setiap 1 detik
    const intervalId = setInterval(updateChart, 1000);

    return () => clearInterval(intervalId);
  }, []); // Dependency array kosong = Timer jalan terus selamanya tanpa gangguan

  return (
    // FIX: Menggunakan h-87.5 (350px / 4) sesuai saran Tailwind
    <div className="bg-slate-900/50 border border-slate-800 rounded-3xl p-6 h-88 shadow-lg backdrop-blur-sm flex flex-col">
      <div className="flex justify-between items-center mb-4 shrink-0">
        <div>
          <h3 className="font-bold text-lg text-white flex items-center gap-2">
            <Activity className="text-blue-500" size={20} />
            Live Traffic Analysis
          </h3>
          <p className="text-slate-400 text-xs mt-1">
            Real-time throughput (Last 20s)
          </p>
        </div>

        <div className="flex gap-4">
          <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20">
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            <span className="text-xs font-bold text-emerald-500">Allowed</span>
          </div>
          <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-red-500/10 border border-red-500/20">
            <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse"></span>
            <span className="text-xs font-bold text-red-500">Blocked</span>
          </div>
        </div>
      </div>

      <div className="flex-1 w-full min-h-0">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={dataHistory}>
            <defs>
              <linearGradient id="colorValid" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="colorBlocked" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
              </linearGradient>
            </defs>
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
              interval="preserveStartEnd"
              minTickGap={20}
            />
            <YAxis
              stroke="#64748b"
              fontSize={10}
              tickLine={false}
              axisLine={false}
              allowDecimals={false}
              width={30}
            />
            <Tooltip
              content={<CustomTooltip />}
              cursor={{
                stroke: "#475569",
                strokeWidth: 1,
                strokeDasharray: "4 4",
              }}
            />
            <Area
              type="monotone"
              dataKey="valid"
              stroke="#10b981"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#colorValid)"
              animationDuration={300}
              isAnimationActive={true}
            />
            <Area
              type="monotone"
              dataKey="blocked"
              stroke="#ef4444"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#colorBlocked)"
              animationDuration={300}
              isAnimationActive={true}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
