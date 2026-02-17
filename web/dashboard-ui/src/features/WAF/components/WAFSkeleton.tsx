import Skeleton from "../../../components/ui/Skeleton";

export default function WAFSkeleton() {
  return (
    <div className="space-y-8">
      {/* Top 3 Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-24 w-full" />
        ))}
      </div>

      {/* Grid: Chart + Table */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <Skeleton className="h-80 lg:col-span-1" />
        <Skeleton className="h-80 lg:col-span-2" />
      </div>
    </div>
  );
}
