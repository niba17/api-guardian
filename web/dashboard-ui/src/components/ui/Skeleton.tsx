interface SkeletonProps {
  className?: string;
}

export default function Skeleton({ className }: SkeletonProps) {
  // animate-pulse adalah class bawaan Tailwind untuk efek berdenyut
  return (
    <div className={`bg-slate-800 animate-pulse rounded-xl ${className}`} />
  );
}
