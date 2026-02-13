import { useMemo } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import type { SecurityLog } from "../../../types/api.types";

interface CountryChartProps {
  logs: SecurityLog[];
}

export default function CountryChart({ logs }: CountryChartProps) {
  // --- LOGIC: Hitung Top 5 Negara Penyerang ---
  const data = useMemo(() => {
    const countryCount: Record<string, number> = {};

    logs.forEach((log) => {
      // Hanya hitung yang DIBLOKIR atau semua? Kita hitung semua traffic dulu.
      const country = log.country || "Unknown";
      countryCount[country] = (countryCount[country] || 0) + 1;
    });

    // Ubah ke Array, Urutkan dari terbanyak, Ambil 5 teratas
    return Object.entries(countryCount)
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 5);
  }, [logs]);

  if (data.length === 0) {
    return (
      <div className="text-slate-500 text-center py-10 font-mono text-sm">
        No data available
      </div>
    );
  }

  return (
    <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6 h-80">
      <h3 className="font-bold text-lg text-white mb-4">
        Top Source Countries
      </h3>

      <ResponsiveContainer width="100%" height="80%">
        <BarChart data={data} layout="vertical" margin={{ left: 20 }}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="#1e293b"
            horizontal={true}
            vertical={false}
          />
          <XAxis type="number" stroke="#64748b" fontSize={12} hide />
          <YAxis
            dataKey="name"
            type="category"
            stroke="#94a3b8"
            fontSize={12}
            width={100}
            tickLine={false}
            axisLine={false}
          />
          <Tooltip
            cursor={{ fill: "#1e293b" }}
            contentStyle={{
              backgroundColor: "#0f172a",
              borderColor: "#334155",
              color: "#fff",
            }}
          />
          <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={20}>
            {/* 👇 Ganti 'entry' jadi '_' */}
            {data.map((_, index) => (
              <Cell
                key={`cell-${index}`}
                fill={index === 0 ? "#f43f5e" : "#3b82f6"}
              />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
