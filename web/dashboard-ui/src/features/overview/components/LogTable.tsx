import {
  CheckCircle,
  Shield,
  Monitor,
  Smartphone,
  Globe,
  Bot,
} from "lucide-react";
import type { SecurityLog } from "../../../types/api.types";
import { formatTime } from "../../../utils/time";

interface LogTableProps {
  logs: SecurityLog[];
  onLogClick: (log: SecurityLog) => void;
}

export default function LogTable({ logs, onLogClick }: LogTableProps) {
  // Helper untuk menentukan icon Device
  const getDeviceIcon = (os?: string, isBot?: boolean) => {
    if (isBot) return <Bot size={14} className="text-orange-400" />;
    const lowerOS = os?.toLowerCase() || "";
    if (lowerOS.includes("win"))
      return <Monitor size={14} className="text-blue-400" />;
    if (lowerOS.includes("ios") || lowerOS.includes("android"))
      return <Smartphone size={14} className="text-purple-400" />;
    return <Globe size={14} className="text-slate-400" />;
  };

  return (
    <div className="bg-guardian-card border border-slate-800 rounded-3xl overflow-hidden shadow-lg">
      <div className="p-6 border-b border-slate-800 flex justify-between items-center">
        <h3 className="font-bold text-lg text-white flex items-center gap-2">
          <ActivityLogIcon /> Live Audit Logs
        </h3>
        <span className="text-xs text-slate-500 font-mono">
          Showing last {logs.length} activities
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead className="bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider">
            <tr>
              <th className="p-4 font-medium">Timestamp</th>
              <th className="p-4 font-medium">Origin & Device</th>
              <th className="p-4 font-medium">Endpoint</th>
              <th className="p-4 font-medium">Response</th>
              <th className="p-4 font-medium">Threat Type</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800 text-sm">
            {logs.length === 0 ? (
              <tr>
                <td
                  colSpan={5}
                  className="p-8 text-center text-slate-500 italic"
                >
                  No logs recorded yet...
                </td>
              </tr>
            ) : (
              logs.map((log) => (
                <tr
                  key={log.id}
                  onClick={() => onLogClick(log)}
                  className="group hover:bg-slate-800/50 transition-colors cursor-pointer"
                >
                  {/* TIMESTAMP */}
                  <td className="p-4 text-slate-400 font-mono whitespace-nowrap">
                    {formatTime(log.timestamp)}
                  </td>

                  {/* SOURCE & DEVICE INTEL */}
                  <td className="p-4">
                    <div className="flex flex-col">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-slate-200 font-mono">
                          {log.ip}
                        </span>
                        <div className="flex items-center gap-1 bg-slate-800 px-1.5 py-0.5 rounded border border-slate-700">
                          {getDeviceIcon(log.os, log.is_bot)}
                          <span className="text-[10px] text-slate-400 uppercase font-bold">
                            {log.browser}
                          </span>
                        </div>
                      </div>
                      <span className="text-xs text-slate-500">
                        {log.city}, {log.country}
                      </span>
                    </div>
                  </td>

                  {/* ENDPOINT */}
                  <td className="p-4">
                    <div className="flex items-center gap-2">
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                          log.method === "GET"
                            ? "bg-blue-500/20 text-blue-400"
                            : log.method === "POST"
                            ? "bg-green-500/20 text-green-400"
                            : log.method === "DELETE"
                            ? "bg-red-500/20 text-red-400"
                            : "bg-slate-700 text-slate-300"
                        }`}
                      >
                        {log.method}
                      </span>
                      <span
                        className="text-slate-300 font-mono truncate max-w-40"
                        title={log.path}
                      >
                        {log.path}
                      </span>
                    </div>
                  </td>

                  {/* PERFORMANCE */}
                  <td className="p-4">
                    <div className="flex flex-col">
                      <span
                        className={`font-mono font-bold ${
                          log.status >= 400
                            ? "text-red-400"
                            : "text-emerald-400"
                        }`}
                      >
                        {log.status}
                      </span>
                      <span className="text-[10px] text-slate-500">
                        {log.latency}ms
                      </span>
                    </div>
                  </td>

                  {/* Threat Type */}
                  <td className="p-2">
                    {log.is_blocked ? (
                      <div className="flex flex-col gap-1">
                        <span className="flex items-center gap-1 text-red-400 text-[11px] font-bold bg-red-400/10 px-2 py-0.5 rounded-full w-fit border border-red-400/20">
                          <Shield size={12} /> {log.threat_type}
                        </span>
                        <span className="text-[10px] text-slate-500 font-mono ml-2">
                          Reason: {log.threat_details}
                        </span>
                      </div>
                    ) : (
                      <span className="flex items-center gap-1 text-emerald-400 text-[11px] font-bold bg-emerald-400/10 px-2 py-0.5 rounded-full w-fit border border-emerald-400/20">
                        <CheckCircle size={12} /> Authorized
                      </span>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function ActivityLogIcon() {
  return (
    <div className="relative flex h-2 w-2">
      <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
      <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
    </div>
  );
}
