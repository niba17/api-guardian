import Skeleton from "../../../components/ui/Skeleton";

export default function TrafficSkeleton() {
  return (
    <div className="space-y-8">
      {/* Top 3 Performance Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-32 w-full" />
        ))}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <Skeleton className="h-80 lg:col-span-2" />
        <Skeleton className="h-80 lg:col-span-1" />
      </div>

      {/* Bottom List Skeleton */}
      <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6 space-y-4">
        <Skeleton className="h-8 w-1/3 mb-4" />
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    </div>
  );
}
