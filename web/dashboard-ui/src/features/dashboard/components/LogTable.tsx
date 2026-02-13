import { CheckCircle, Shield } from "lucide-react"; // 👈 Import hanya yang dipakai
import type { SecurityLog } from "../../../types/api.types";

interface LogTableProps {
  logs: SecurityLog[];
  onLogClick: (log: SecurityLog) => void;
}

export default function LogTable({ logs, onLogClick }: LogTableProps) {
  return (
    <div className="bg-guardian-card border border-slate-800 rounded-3xl overflow-hidden shadow-lg">
      <div className="p-6 border-b border-slate-800 flex justify-between items-center">
        <h3 className="font-bold text-lg text-white flex items-center gap-2">
          <ActivityLogIcon /> Live Audit Logs
        </h3>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead className="bg-slate-900/50 text-slate-400 text-xs uppercase tracking-wider">
            <tr>
              <th className="p-4 font-medium">Timestamp</th>
              <th className="p-4 font-medium">Source IP</th>
              <th className="p-4 font-medium">Endpoint</th>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium">Action</th>
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
                  <td className="p-4 text-slate-400 font-mono whitespace-nowrap">
                    {new Date(log.timestamp).toLocaleTimeString()}
                  </td>
                  <td className="p-4">
                    <div className="flex flex-col">
                      <span className="font-medium text-slate-200 font-mono">
                        {log.ip}
                      </span>
                      <span className="text-xs text-slate-500">
                        {log.city}, {log.country}
                      </span>
                    </div>
                  </td>
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
                      {/* 👇 FIX: Ganti max-w-[200px] jadi max-w-50 */}
                      <span className="text-slate-300 font-mono truncate max-w-50">
                        {log.path}
                      </span>
                    </div>
                  </td>
                  <td className="p-4">
                    <span
                      className={`font-mono font-bold ${
                        log.status >= 400 ? "text-red-400" : "text-emerald-400"
                      }`}
                    >
                      {log.status}
                    </span>
                    <span className="text-xs text-slate-500 ml-2">
                      {log.latency}ms
                    </span>
                  </td>
                  <td className="p-4">
                    {log.is_blocked ? (
                      <span className="flex items-center gap-1 text-red-400 text-xs font-bold bg-red-400/10 px-2 py-1 rounded-full w-fit">
                        <Shield size={12} /> Blocked
                      </span>
                    ) : (
                      <span className="flex items-center gap-1 text-emerald-400 text-xs font-bold bg-emerald-400/10 px-2 py-1 rounded-full w-fit">
                        <CheckCircle size={12} /> Allowed
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

// Komponen Ikon Kedip-kedip (Visual Only)
function ActivityLogIcon() {
  return (
    <div className="relative flex h-2 w-2">
      <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
      <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
    </div>
  );
}
