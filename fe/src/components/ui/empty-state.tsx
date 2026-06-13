import React from 'react';
import Link from 'next/link';
import { LucideIcon } from 'lucide-react';
import { Button } from './button';

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: {
    label: string;
    href: string;
  };
}

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="glass-card p-12 text-center border border-white/[0.06] shadow-xl flex flex-col items-center justify-center max-w-lg mx-auto my-8">
      <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-accent/10 mb-6 animate-pulse-glow">
        <Icon className="w-8 h-8 text-accent" />
      </div>
      <h3 className="text-xl font-semibold text-slate-100 mb-2">{title}</h3>
      <p className="text-slate-400 mb-6 max-w-sm">{description}</p>
      {action && (
        <Link href={action.href}>
          <Button size="md">{action.label}</Button>
        </Link>
      )}
    </div>
  );
}
