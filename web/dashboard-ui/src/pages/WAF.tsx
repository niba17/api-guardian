import { useState } from "react";
import MainLayout from "../components/layout/MainLayout";
import {
  AlertOctagon,
  Map as MapIcon, // Rename import biar gak bentrok sama Map Constructor
  Target,
  Activity,
  ShieldX,
  Terminal,
} from "lucide-react";
import { useWAFData } from "../features/WAF/hooks/useWAF"; // 👈 Hook Baru
import CountryChart from "../features/WAF/components/CountryChart";
import LogTable from "../features/overview/components/LogTable";
import ThreatsSkeleton from "../features/WAF/components/WAFSkeleton";
import Modal from "../components/ui/Modal";
import type { SecurityLog } from "../types/api.types";

export default function WAF() {
  // Ambil data yang sudah matang dari hook WAF
  const { loading, threats, totalThreats, uniqueCountries, topTarget } =
    useWAFData();
  const [selectedLog, setSelectedLog] = useState<SecurityLog | null>(null);

  return (
    <MainLayout>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
          WAF Monitor
          {loading && (
            <Activity
              className="animate-spin text-guardian-warning"
              size={20}
            />
          )}
        </h1>
        <p className="text-slate-400">
          Real-time analysis of blocked requests and potential security
          breaches.
        </p>
      </div>

      {loading ? (
        <ThreatsSkeleton />
      ) : (
        <>
          {/* KPI Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <StatCard
              title="Blocked Threats"
              value={totalThreats}
              icon={<AlertOctagon size={24} />}
              color="red"
            />
            <StatCard
              title="Source Countries"
              value={uniqueCountries}
              icon={<MapIcon size={24} />}
              color="blue"
            />
            <StatCard
              title="Active Target"
              value={topTarget}
              icon={<Target size={24} />}
              color="yellow"
              isText
            />
          </div>

          {/* Main Content */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1">
              {/* Kirim data threats (yang sudah difilter blocked) ke Chart */}
              <CountryChart logs={threats} />
            </div>
            <div className="lg:col-span-2">
              <LogTable logs={threats} onLogClick={setSelectedLog} />
            </div>
          </div>
        </>
      )}

      {/* Forensic Modal */}
      <Modal
        isOpen={!!selectedLog}
        onClose={() => setSelectedLog(null)}
        title="Threat Forensic Detail"
      >
        {selectedLog && <ForensicDetail log={selectedLog} />}
      </Modal>
    </MainLayout>
  );
}

// --- SUB-COMPONENTS (Biar kode utama bersih) ---

// --- DEFINISI TIPE DATA (Fix ESLint no-explicit-any) ---

interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode; // Tipe untuk elemen Icon/JSX
  color: string;
  isText?: boolean; // Tanda tanya (?) artinya optional
}

function StatCard({ title, value, icon, color, isText }: StatCardProps) {
  return (
    <div
      className={`bg-${color}-500/10 border border-${color}-500/20 p-6 rounded-2xl flex items-center gap-4`}
    >
      <div className={`p-3 bg-${color}-500/20 rounded-lg text-${color}-500`}>
        {icon}
      </div>
      <div>
        <p className="text-slate-500 text-[10px] font-bold uppercase tracking-widest">
          {title}
        </p>
        <h3
          className={`font-bold text-white font-mono ${
            isText ? "text-xl" : "text-3xl"
          }`}
        >
          {value}
        </h3>
      </div>
    </div>
  );
}

function ForensicDetail({ log }: { log: SecurityLog }) {
  return (
    <div className="space-y-6">
      <div className="p-4 rounded-xl border flex items-center gap-4 bg-red-500/10 border-red-500/30 text-red-400">
        <ShieldX size={32} />
        <div>
          <h4 className="font-bold text-lg">Intrusion Blocked</h4>
          <p className="text-sm opacity-80">
            Origin IP: <span className="font-mono font-bold">{log.ip}</span>
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
          {JSON.stringify(log, null, 2)}
        </pre>
      </div>
    </div>
  );
}
