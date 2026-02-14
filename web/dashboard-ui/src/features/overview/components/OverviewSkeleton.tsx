import Skeleton from "../../../components/ui/Skeleton";

export default function OverviewSkeleton() {
  return (
    <div className="space-y-8">
      {/* 1. Stats Cards Skeleton (5 Kolom) */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-32 w-full" />
        ))}
      </div>

      {/* 2. Chart Skeleton */}
      <Skeleton className="h-87.5 w-full" />

      {/* 3. Table Skeleton */}
      <div className="bg-guardian-card border border-slate-800 rounded-3xl p-6 space-y-4">
        <Skeleton className="h-8 w-1/4" /> {/* Judul Tabel */}
        <div className="space-y-2">
          {[...Array(6)].map((_, i) => (
            <Skeleton key={i} className="h-12 w-full" />
          ))}
        </div>
      </div>
    </div>
  );
}
