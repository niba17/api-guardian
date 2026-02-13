import { useMemo } from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";
import type { SecurityLog } from "../../../types/api.types";

interface MethodChartProps {
  logs: SecurityLog[];
}

const COLORS = ["#3b82f6", "#10b981", "#f59e0b", "#ef4444"]; // Biru, Hijau, Kuning, Merah

export default function MethodChart({ logs }: MethodChartProps) {
  const data = useMemo(() => {
    const methodCount: Record<string, number> = {};

    logs.forEach((log) => {
      const method = log.method || "UNKNOWN";
      methodCount[method] = (methodCount[method] || 0) + 1;
    });

    return Object.entries(methodCount).map(([name, value]) => ({
      name,
      value,
    }));
  }, [logs]);

  if (data.length === 0)
    return <div className="text-slate-500 text-center py-10">No data</div>;

  return (
    <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6 h-80">
      <h3 className="font-bold text-lg text-white mb-4">HTTP Methods</h3>
      <ResponsiveContainer width="100%" height="85%">
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            innerRadius={60} // Biar jadi Donut Chart (Bolong tengah)
            outerRadius={80}
            paddingAngle={5}
            dataKey="value"
          >
            {data.map((_, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
          <Tooltip
            contentStyle={{
              backgroundColor: "#0f172a",
              borderColor: "#334155",
              borderRadius: "8px",
            }}
            itemStyle={{ color: "#fff" }}
          />
          <Legend verticalAlign="bottom" height={36} iconType="circle" />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
