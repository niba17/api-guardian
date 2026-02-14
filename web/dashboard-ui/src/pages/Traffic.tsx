import MainLayout from "../components/layout/MainLayout";
import { Activity, Server, ArrowUpRight, CloudLightning } from "lucide-react";
import { useTrafficData } from "../features/traffic/hooks/useTraffic"; // 👈 Pakai hook baru
import MethodChart from "../features/traffic/components/MethodChart";
import LatencyChart from "../features/traffic/components/LatencyChart";
import TrafficSkeleton from "../features/traffic/components/TrafficSkeleton";

export default function Traffic() {
  // Ambil data yang sudah "matang" dari hook
  const { logs, stats, loading, topEndpoints } = useTrafficData();

  return (
    <MainLayout>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
            Traffic
            {loading && (
              <Activity
                className="animate-spin text-guardian-warning"
                size={20}
              />
            )}
          </h1>
          <p className="text-slate-400">
            Deep packet inspection and performance metrics.
          </p>
        </div>
      </div>

      {loading ? (
        <TrafficSkeleton />
      ) : (
        <>
          {/* KPI Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <PerformanceCard
              title="Avg Response Time"
              value={stats?.avg_latency || "0ms"}
              icon={<Server />}
              color="blue"
              healthTag
            />
            <PerformanceCard
              title="Total Request Processed"
              value={stats?.total_requests || 0}
              icon={<CloudLightning />}
              color="purple"
            />
            <PerformanceCard
              title="Success Rate"
              value={`${
                stats?.total_requests
                  ? Math.round(
                      ((stats.total_requests - stats.blocked_requests) /
                        stats.total_requests) *
                        100
                    )
                  : 100
              }%`}
              icon={<ArrowUpRight />}
              color="emerald"
            />
          </div>

          {/* Charts */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
            <div className="lg:col-span-2">
              <LatencyChart logs={logs} />
            </div>
            <div className="lg:col-span-1">
              <MethodChart logs={logs} />
            </div>
          </div>

          {/* Top Endpoints List */}
          <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6">
            <h3 className="font-bold text-lg text-white mb-6">
              Most Accessed Endpoints
            </h3>
            <div className="space-y-4">
              {topEndpoints.length === 0 ? (
                <p className="text-slate-500 text-sm">No data yet...</p>
              ) : (
                topEndpoints.map(([path, count], index) => (
                  <EndpointRow
                    key={path}
                    path={path}
                    count={count}
                    index={index}
                    maxCount={topEndpoints[0][1]}
                  />
                ))
              )}
            </div>
          </div>
        </>
      )}
    </MainLayout>
  );
}

// --- DEFINISI TIPE DATA (Fix ESLint no-explicit-any) ---

interface PerformanceCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode; // Tipe untuk Icon Component
  color: string;
  healthTag?: boolean; // Tanda tanya (?) artinya optional
}

function PerformanceCard({
  title,
  value,
  icon,
  color,
  healthTag,
}: PerformanceCardProps) {
  return (
    <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
      <div className="flex justify-between items-start mb-4">
        {/* Note: Tailwind dynamic classes (bg-${color}) kadang tricky, 
            pastikan color yang dikirim valid (blue, purple, emerald) */}
        <div className={`p-2 bg-${color}-500/10 rounded-lg text-${color}-400`}>
          {icon}
        </div>
        {healthTag && (
          <span className="text-xs font-bold bg-green-500/10 text-green-500 px-2 py-1 rounded">
            HEALTHY
          </span>
        )}
      </div>
      <p className="text-slate-500 text-xs uppercase font-bold">{title}</p>
      <h3 className="text-3xl font-bold text-white">{value}</h3>
    </div>
  );
}

interface EndpointRowProps {
  path: string;
  count: number;
  index: number;
  maxCount: number;
}

function EndpointRow({ path, count, index, maxCount }: EndpointRowProps) {
  return (
    <div className="flex items-center justify-between p-3 bg-slate-900/50 rounded-xl border border-slate-800">
      <div className="flex items-center gap-4">
        <span className="flex items-center justify-center w-6 h-6 rounded-full bg-slate-800 text-xs font-bold text-slate-400">
          {index + 1}
        </span>
        <span className="text-blue-400 font-mono text-sm">{path}</span>
      </div>
      <div className="flex items-center gap-2">
        <div className="h-2 w-24 bg-slate-800 rounded-full overflow-hidden">
          <div
            className="h-full bg-blue-500"
            style={{ width: `${(count / maxCount) * 100}%` }}
          ></div>
        </div>
        <span className="text-white font-bold text-sm">{count}</span>
      </div>
    </div>
  );
}
