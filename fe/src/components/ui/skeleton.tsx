import React from 'react';
import { cn } from '@/lib/utils';

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {}

export function Skeleton({ className, ...props }: SkeletonProps) {
  return (
    <div
      className={cn('animate-pulse rounded bg-white/[0.06]', className)}
      {...props}
    />
  );
}

export function SkeletonCard() {
  return (
    <div className="glass-card p-4 flex flex-col h-full border border-white/[0.06] space-y-4">
      {/* Product Image Area */}
      <Skeleton className="w-full aspect-square rounded-lg" />
      {/* Product Name */}
      <Skeleton className="h-6 w-3/4 rounded" />
      {/* Product Description */}
      <div className="space-y-2 flex-grow">
        <Skeleton className="h-4 w-full rounded" />
        <Skeleton className="h-4 w-5/6 rounded" />
      </div>
      {/* Footer Price & Add To Cart Button */}
      <div className="flex items-center justify-between pt-2">
        <Skeleton className="h-6 w-1/4 rounded" />
        <Skeleton className="h-9 w-24 rounded-lg" />
      </div>
    </div>
  );
}

export function SkeletonTableRow() {
  return (
    <tr className="animate-pulse border-b border-white/[0.06]">
      <td className="px-6 py-4">
        <Skeleton className="h-5 w-16 rounded" />
      </td>
      <td className="px-6 py-4">
        <Skeleton className="h-5 w-24 rounded" />
      </td>
      <td className="px-6 py-4">
        <Skeleton className="h-6 w-20 rounded-full" />
      </td>
      <td className="px-6 py-4">
        <Skeleton className="h-5 w-32 rounded" />
      </td>
      <td className="px-6 py-4">
        <Skeleton className="h-8 w-16 rounded" />
      </td>
    </tr>
  );
}
