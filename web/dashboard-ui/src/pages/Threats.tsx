import { useState } from "react";
import MainLayout from "../components/layout/MainLayout";
import {
  AlertOctagon,
  Map,
  Target,
  Activity,
  ShieldX, // 👈 Tambah ini
  Terminal, // 👈 Tambah ini
} from "lucide-react";
import { useDashboardData } from "../features/dashboard/api/useDashboard";
import CountryChart from "../features/dashboard/components/CountryChart";
import LogTable from "../features/dashboard/components/LogTable";
import ThreatsSkeleton from "../features/dashboard/components/ThreatsSkeleton";
import Modal from "../components/ui/Modal"; // 👈 Tambah ini
import type { SecurityLog } from "../types/api.types";

export default function Threats() {
  const { logs, loading } = useDashboardData();
  const [selectedLog, setSelectedLog] = useState<SecurityLog | null>(null);

  const threats = logs.filter((log) => log.is_blocked);
  const totalThreats = threats.length;
  const uniqueCountries = new Set(threats.map((l) => l.country)).size;

  return (
    <MainLayout>
      {/* 1. Header - Flat & Clean */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
          {/* <ShieldAlert className="text-red-500" /> */}
          Threat Intelligence
          {loading && (
            <Activity
              className="animate-spin text-guardian-warning"
              size={20}
            />
          )}
        </h1>
        <p className="text-slate-400">
          Analysis of blocked requests and potential security breaches.
        </p>
      </div>

      {loading ? (
        <ThreatsSkeleton />
      ) : (
        <>
          {/* 2. Grid Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-red-500/10 border border-red-500/20 p-6 rounded-2xl flex items-center gap-4">
              <div className="p-3 bg-red-500/20 rounded-lg text-red-500">
                <AlertOctagon size={24} />
              </div>
              <div>
                <p className="text-slate-500 text-[10px] font-bold uppercase tracking-widest">
                  Blocked Threats
                </p>
                <h3 className="text-3xl font-bold text-white font-mono">
                  {totalThreats}
                </h3>
              </div>
            </div>

            <div className="bg-blue-500/10 border border-blue-500/20 p-6 rounded-2xl flex items-center gap-4">
              <div className="p-3 bg-blue-500/20 rounded-lg text-blue-500">
                <Map size={24} />
              </div>
              <div>
                <p className="text-slate-500 text-[10px] font-bold uppercase tracking-widest">
                  Source Countries
                </p>
                <h3 className="text-3xl font-bold text-white font-mono">
                  {uniqueCountries}
                </h3>
              </div>
            </div>

            <div className="bg-yellow-500/10 border border-yellow-500/20 p-6 rounded-2xl flex items-center gap-4">
              <div className="p-3 bg-yellow-500/20 rounded-lg text-yellow-500">
                <Target size={24} />
              </div>
              <div>
                <p className="text-slate-500 text-[10px] font-bold uppercase tracking-widest">
                  Active Target
                </p>
                <h3 className="text-xl font-bold text-white font-mono">
                  /api/v1/auth
                </h3>
              </div>
            </div>
          </div>

          {/* 3. Main Content Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1">
              <CountryChart logs={threats} />
            </div>
            <div className="lg:col-span-2">
              <LogTable logs={threats} onLogClick={setSelectedLog} />
            </div>
          </div>
        </>
      )}

      {/* 4. Modal Detail - FIX: Agar selectedLog terbaca */}
      <Modal
        isOpen={!!selectedLog}
        onClose={() => setSelectedLog(null)}
        title="Threat Detail"
      >
        {selectedLog && (
          <div className="space-y-6">
            <div className="p-4 rounded-xl border flex items-center gap-4 bg-red-500/10 border-red-500/30 text-red-400">
              <ShieldX size={32} />
              <div>
                <h4 className="font-bold text-lg">Intrusion Blocked</h4>
                <p className="text-sm opacity-80">
                  Origin IP:{" "}
                  <span className="font-mono font-bold">{selectedLog.ip}</span>
                </p>
              </div>
            </div>

            <div className="bg-black/80 rounded-lg p-4 border border-slate-800 overflow-x-auto">
              <div className="flex items-center gap-2 mb-2 text-slate-400 border-b border-slate-800 pb-2">
                <Terminal size={16} />
                <span className="text-xs font-bold uppercase tracking-tighter">
                  Forensic Raw Evidence
                </span>
              </div>
              <pre className="text-xs font-mono text-red-400">
                {JSON.stringify(selectedLog, null, 2)}
              </pre>
            </div>
          </div>
        )}
      </Modal>
    </MainLayout>
  );
}
