import Skeleton from "../../../components/ui/Skeleton";

export default function ConfigSkeleton() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
      {[...Array(2)].map((_, i) => (
        <div
          key={i}
          className="bg-guardian-card border border-slate-800 rounded-2xl p-6 space-y-6"
        >
          <Skeleton className="h-8 w-1/2 mb-4" /> {/* Header */}
          <div className="space-y-8">
            {[...Array(3)].map((_, j) => (
              <div key={j} className="flex justify-between items-center">
                <div className="space-y-2 w-2/3">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-3 w-1/2" />
                </div>
                <Skeleton className="h-6 w-12 rounded-full" /> {/* Switch */}
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
