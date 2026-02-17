import type { LucideIcon } from "lucide-react";

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  color: string;
}

export default function StatsCard({
  title,
  value,
  icon: Icon,
  color,
}: StatsCardProps) {
  return (
    <div className="bg-guardian-card p-6 rounded-3xl border border-slate-800 hover:border-slate-700 transition-all cursor-default group">
      <div className="flex justify-between items-start mb-4">
        <div
          className={`p-3 rounded-2xl bg-slate-900/50 ${color} group-hover:scale-110 transition-transform duration-300`}
        >
          <Icon size={24} />
        </div>
      </div>
      <p className="text-slate-400 text-sm font-medium">{title}</p>
      <p className="text-3xl font-bold text-white mt-1 font-mono tracking-tight">
        {value}
      </p>
    </div>
  );
}
