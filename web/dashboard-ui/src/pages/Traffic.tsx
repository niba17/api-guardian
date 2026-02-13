import MainLayout from "../components/layout/MainLayout";
import { Activity, Server, ArrowUpRight, CloudLightning } from "lucide-react";
import { useDashboardData } from "../features/dashboard/api/useDashboard";
import MethodChart from "../features/dashboard/components/MethodChart";
import LatencyChart from "../features/dashboard/components/LatencyChart";
import TrafficSkeleton from "../features/dashboard/components/TrafficSkeleton";

export default function Traffic() {
  const { logs, stats, loading } = useDashboardData();

  // --- LOGIC: Hitung Top Endpoint ---
  const topEndpoints = logs.reduce((acc, log) => {
    acc[log.path] = (acc[log.path] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const sortedEndpoints = Object.entries(topEndpoints)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 5);

  return (
    <MainLayout>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
            {/* <Activity className="text-blue-400" /> */}
            Traffic Monitor
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

      {/* --- SWITCH: SKELETON VS CONTENT --- */}
      {loading ? (
        <TrafficSkeleton />
      ) : (
        <>
          {/* 1. KPI Cards (Performance Focus) */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
              <div className="flex justify-between items-start mb-4">
                <div className="p-2 bg-blue-500/10 rounded-lg text-blue-400">
                  <Server />
                </div>
                <span className="text-xs font-bold bg-green-500/10 text-green-500 px-2 py-1 rounded">
                  HEALTHY
                </span>
              </div>
              <p className="text-slate-500 text-xs uppercase font-bold">
                Avg Response Time
              </p>
              <h3 className="text-3xl font-bold text-white">
                {stats?.avg_latency || "0ms"}
              </h3>
            </div>

            <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
              <div className="flex justify-between items-start mb-4">
                <div className="p-2 bg-purple-500/10 rounded-lg text-purple-400">
                  <CloudLightning />
                </div>
              </div>
              <p className="text-slate-500 text-xs uppercase font-bold">
                Total Request Processed
              </p>
              <h3 className="text-3xl font-bold text-white">
                {stats?.total_requests || 0}
              </h3>
            </div>

            <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl">
              <div className="flex justify-between items-start mb-4">
                <div className="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
                  <ArrowUpRight />
                </div>
              </div>
              <p className="text-slate-500 text-xs uppercase font-bold">
                Success Rate
              </p>
              <h3 className="text-3xl font-bold text-white">
                {stats?.total_requests
                  ? Math.round(
                      ((stats.total_requests - stats.blocked_requests) /
                        stats.total_requests) *
                        100
                    )
                  : 100}
                %
              </h3>
            </div>
          </div>

          {/* 2. Charts Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
            <div className="lg:col-span-2">
              <LatencyChart logs={logs} />
            </div>
            <div className="lg:col-span-1">
              <MethodChart logs={logs} />
            </div>
          </div>

          {/* 3. Top Endpoints List */}
          <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6">
            <h3 className="font-bold text-lg text-white mb-6">
              Most Accessed Endpoints
            </h3>
            <div className="space-y-4">
              {sortedEndpoints.length === 0 ? (
                <p className="text-slate-500 text-sm">No data yet...</p>
              ) : (
                sortedEndpoints.map(([path, count], index) => (
                  <div
                    key={path}
                    className="flex items-center justify-between p-3 bg-slate-900/50 rounded-xl border border-slate-800"
                  >
                    <div className="flex items-center gap-4">
                      <span className="flex items-center justify-center w-6 h-6 rounded-full bg-slate-800 text-xs font-bold text-slate-400">
                        {index + 1}
                      </span>
                      <span className="text-blue-400 font-mono text-sm">
                        {path}
                      </span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="h-2 w-24 bg-slate-800 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-blue-500"
                          style={{
                            width: `${(count / sortedEndpoints[0][1]) * 100}%`,
                          }}
                        ></div>
                      </div>
                      <span className="text-white font-bold text-sm">
                        {count}
                      </span>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </>
      )}
    </MainLayout>
  );
}
