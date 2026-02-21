import { useState } from "react";
import MainLayout from "../components/layout/MainLayout";
import StatsCard from "../features/overview/components/StatsCard";
import TrafficChart from "../features/traffic/components/TrafficChart";
import LogTable from "../features/overview/components/LogTable";
import Modal from "../components/ui/Modal";
import DashboardSkeleton from "../features/overview/components/OverviewSkeleton";
import {
  Zap,
  ShieldX,
  Globe,
  Cpu,
  Activity,
  CheckCircle,
  Terminal,
  Monitor,
  Bot,
  User,
  Eye,
} from "lucide-react";
import { useOverviewData } from "../features/overview/hooks/useOverview";
import type { SecurityLog } from "../types/api.types";
import HistoryChart from "../features/overview/components/HistoryChart";

export default function Overview() {
  const { stats, logs, loading, lastUpdate } = useOverviewData();
  const [selectedLog, setSelectedLog] = useState<SecurityLog | null>(null);

  // Kalkulasi data valid
  const total = stats?.total_requests || 0;
  const blocked = stats?.total_blocked || 0;
  const valid = total - blocked;

  return (
    <MainLayout>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
            Overview
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
            <HistoryChart logs={logs} />
          </div>

          <div className="grid grid-cols-1 gap-8">
            <LogTable logs={logs} onLogClick={setSelectedLog} />
          </div>
        </>
      )}

      {/* --- MODAL FORENSIK TERPADU --- */}
      <Modal
        isOpen={!!selectedLog}
        onClose={() => setSelectedLog(null)}
        title="Forensic Investigation"
      >
        {selectedLog && (
          <div className="space-y-6">
            {/* Header Status */}
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
                    ? `Threat Blocked: ${
                        selectedLog.threat_type || "WAF Alert"
                      }`
                    : "Request Authorized"}
                </h4>
                <p className="text-sm opacity-80">
                  HTTP {selectedLog.status} • {selectedLog.latency}ms Latency
                </p>
              </div>
            </div>

            {/* Source & Target Grid */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-800">
                <p className="text-xs text-slate-500 uppercase font-bold mb-1">
                  Origin IP
                </p>
                <p className="text-white font-mono text-lg">{selectedLog.ip}</p>
                <p className="text-sm text-slate-400">
                  {selectedLog.city}, {selectedLog.country}
                </p>
              </div>
              <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-800">
                <p className="text-xs text-slate-500 uppercase font-bold mb-1">
                  Target Endpoint
                </p>
                <div className="flex items-center gap-2 mt-1 overflow-hidden">
                  <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-700 text-white">
                    {selectedLog.method}
                  </span>
                  <p className="text-white font-mono text-sm truncate">
                    {selectedLog.path}
                  </p>
                </div>
              </div>
            </div>

            {/* Device Intelligence Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-slate-900/30 p-3 rounded-lg border border-slate-800/50 flex items-center gap-3">
                <div className="p-2 bg-blue-500/10 rounded-lg text-blue-400">
                  <Monitor size={18} />
                </div>
                <div>
                  <p className="text-[10px] text-slate-500 uppercase font-bold">
                    OS
                  </p>
                  <p className="text-sm text-white font-medium">
                    {selectedLog.os || "Unknown"}
                  </p>
                </div>
              </div>
              <div className="bg-slate-900/30 p-3 rounded-lg border border-slate-800/50 flex items-center gap-3">
                <div className="p-2 bg-purple-500/10 rounded-lg text-purple-400">
                  <Globe size={18} />
                </div>
                <div>
                  <p className="text-[10px] text-slate-500 uppercase font-bold">
                    Browser
                  </p>
                  <p className="text-sm text-white font-medium">
                    {selectedLog.browser || "Unknown"}
                  </p>
                </div>
              </div>
              <div className="bg-slate-900/30 p-3 rounded-lg border border-slate-800/50 flex items-center gap-3">
                <div
                  className={`p-2 rounded-lg ${
                    selectedLog.is_bot
                      ? "bg-orange-500/10 text-orange-400"
                      : "bg-emerald-500/10 text-emerald-400"
                  }`}
                >
                  {selectedLog.is_bot ? <Bot size={18} /> : <User size={18} />}
                </div>
                <div>
                  <p className="text-[10px] text-slate-500 uppercase font-bold">
                    Client Type
                  </p>
                  <p className="text-sm text-white font-medium">
                    {selectedLog.is_bot ? "Automated Bot" : "Human User"}
                  </p>
                </div>
              </div>
            </div>

            {/* Intercepted Payload (Body) */}
            {selectedLog.body && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-guardian-warning">
                  <Eye size={16} />
                  <span className="text-xs font-bold uppercase tracking-wider">
                    Intercepted Payload
                  </span>
                </div>
                <div className="bg-black/40 rounded-lg p-4 border border-guardian-warning/20 relative group">
                  <pre className="text-xs font-mono text-guardian-warning/90 whitespace-pre-wrap break-all leading-relaxed">
                    {selectedLog.body}
                  </pre>
                  <div className="absolute top-2 right-2 px-2 py-0.5 bg-guardian-warning/10 rounded text-[8px] text-guardian-warning font-bold border border-guardian-warning/20">
                    PII REDACTED
                  </div>
                </div>
              </div>
            )}

            {/* Technical Metadata (JSON) */}
            <div>
              <div className="flex items-center gap-2 mb-2 text-slate-500">
                <Terminal size={14} />
                <span className="text-[10px] font-bold uppercase tracking-widest">
                  Full Security Metadata
                </span>
              </div>
              <div className="bg-black/60 rounded-lg p-4 overflow-x-auto border border-slate-800/50 max-h-48 overflow-y-auto">
                <pre className="text-[10px] font-mono text-slate-500 leading-tight">
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
