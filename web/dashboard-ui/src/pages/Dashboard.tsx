import { useState } from "react";
import MainLayout from "../components/layout/MainLayout";
import StatsCard from "../features/dashboard/components/StatsCard";
import TrafficChart from "../features/dashboard/components/TrafficChart";
import LogTable from "../features/dashboard/components/LogTable";
import Modal from "../components/ui/Modal";
import DashboardSkeleton from "../features/dashboard/components/DashboardSkeleton"; // 👈 Import Skeleton Baru
import {
  Zap,
  ShieldX,
  Globe,
  Cpu,
  Activity,
  CheckCircle,
  Terminal,
} from "lucide-react";
import { useDashboardData } from "../features/dashboard/api/useDashboard";
import type { SecurityLog } from "../types/api.types";
import HistoryChart from "../features/dashboard/components/HistoryChart";

export default function Dashboard() {
  const { stats, logs, loading, lastUpdate } = useDashboardData();
  const [selectedLog, setSelectedLog] = useState<SecurityLog | null>(null);

  // Kalkulasi data valid
  const total = stats?.total_requests || 0;
  const blocked = stats?.blocked_requests || 0;
  const valid = total - blocked;

  return (
    <MainLayout>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
            Network Overview
            {/* Spinner kecil tetap ada sebagai indikator refresh data di background */}
            {loading && (
              <Activity
                className="animate-spin text-guardian-warning"
                size={20}
              />
            )}
          </h1>
          <p className="text-slate-400">
            Monitoring real-time traffic from API Guardian.
            <span className="text-xs ml-2 opacity-50">
              {loading ? "Synchronizing..." : `Last Update: ${lastUpdate}`}
            </span>
          </p>
        </div>

        <div className="flex gap-2">
          <span
            className={`px-3 py-1 rounded-full text-xs font-mono border flex items-center gap-2 ${
              stats
                ? "bg-guardian-primary/10 text-guardian-primary border-guardian-primary/20"
                : "bg-red-500/10 text-red-500 border-red-500/20"
            }`}
          >
            {stats ? "● SYSTEM ONLINE" : "○ CONNECTING..."}
          </span>
        </div>
      </div>

      {/* --- LOGIC UTAMA: SKELETON VS CONTENT --- */}
      {loading ? (
        <DashboardSkeleton />
      ) : (
        <>
          {/* Grid Stats */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6 mb-8">
            <StatsCard
              title="Total Traffic"
              value={total}
              icon={Zap}
              color="text-blue-400"
            />
            <StatsCard
              title="Valid Requests"
              value={valid}
              icon={CheckCircle}
              color="text-emerald-400"
            />
            <StatsCard
              title="Blocked Attacks"
              value={blocked}
              icon={ShieldX}
              color="text-guardian-danger"
            />
            <StatsCard
              title="Unique IPs"
              value={stats?.unique_ips || 0}
              icon={Globe}
              color="text-guardian-primary"
            />
            <StatsCard
              title="Avg Latency"
              value={stats?.avg_latency || "0ms"}
              icon={Cpu}
              color="text-yellow-400"
            />
          </div>

          <div className="mb-8">
            <TrafficChart logs={logs} />
          </div>

          <div className="mb-8">
            <HistoryChart logs={logs} /> {/* 👈 GANTI INI */}
          </div>

          {/* Baris 3: Tabel Detail (Optional, kalau mau tetap ada) */}
          <div className="grid grid-cols-1 gap-8">
            <LogTable logs={logs} onLogClick={setSelectedLog} />
          </div>
        </>
      )}

      {/* Modal Forensik tetap di luar loading agar tidak ter-reset saat data refresh */}
      <Modal
        isOpen={!!selectedLog}
        onClose={() => setSelectedLog(null)}
        title="Forensic Detail"
      >
        {selectedLog && (
          <div className="space-y-6">
            <div
              className={`p-4 rounded-xl border flex items-center gap-4 ${
                selectedLog.is_blocked
                  ? "bg-red-500/10 border-red-500/30 text-red-400"
                  : "bg-emerald-500/10 border-emerald-500/30 text-emerald-400"
              }`}
            >
              {selectedLog.is_blocked ? (
                <ShieldX size={32} />
              ) : (
                <CheckCircle size={32} />
              )}
              <div>
                <h4 className="font-bold text-lg">
                  {selectedLog.is_blocked
                    ? "Threat Blocked"
                    : "Request Allowed"}
                </h4>
                <p className="text-sm opacity-80">
                  Status Code:{" "}
                  <span className="font-mono font-bold">
                    {selectedLog.status}
                  </span>{" "}
                  • Latency:{" "}
                  <span className="font-mono">{selectedLog.latency}ms</span>
                </p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-800">
                <p className="text-xs text-slate-500 uppercase font-bold mb-1">
                  Source
                </p>
                <p className="text-white font-mono text-lg">{selectedLog.ip}</p>
                <p className="text-sm text-slate-400">
                  {selectedLog.city}, {selectedLog.country}
                </p>
              </div>
              <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-800">
                <p className="text-xs text-slate-500 uppercase font-bold mb-1">
                  Target
                </p>
                <div className="flex items-center gap-2 mt-1">
                  <span className="px-2 py-0.5 rounded text-xs font-bold bg-slate-700 text-white">
                    {selectedLog.method}
                  </span>
                  <p className="text-white font-mono truncate">
                    {selectedLog.path}
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-center gap-2 mb-2 text-slate-400">
                <Terminal size={16} />
                <span className="text-xs font-bold uppercase">
                  Raw JSON Data
                </span>
              </div>
              <div className="bg-black/80 rounded-lg p-4 overflow-x-auto border border-slate-800">
                <pre className="text-xs font-mono text-green-400">
                  {JSON.stringify(selectedLog, null, 2)}
                </pre>
              </div>
            </div>
          </div>
        )}
      </Modal>
    </MainLayout>
  );
}
