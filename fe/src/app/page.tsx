import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { config } from '@/lib/config';
import { Zap, ArrowRight, Cpu, ShieldCheck, ShoppingBag } from 'lucide-react';
import FeaturedProducts from './_components/featured-products';

export default function HomePage() {
  return (
    <div>
      {/* Hero Section */}
      <section className="relative overflow-hidden border-b border-white/[0.06]">
        {/* Background effects */}
        <div className="absolute inset-0 bg-hero-glow" />
        <div className="absolute inset-0 bg-grid opacity-30" />
        
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32">
          <div className="text-center">
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-accent/10 border border-accent/20 text-accent text-sm font-medium mb-8 animate-fade-in">
              <Zap className="w-4 h-4" />
              Discover Next-Gen Commerce
            </div>
            
            <h1 className="text-4xl md:text-6xl font-bold mb-6 text-gradient tracking-tight leading-tight">
              Welcome to {config.appName}
            </h1>
            <p className="text-lg md:text-xl text-slate-400 mb-8 max-w-2xl mx-auto">
              {config.appDescription}. Experience a seamless, ultra-fast purchase flow driven by modern, high-performance cloud intelligence.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link href="/products">
                <Button size="lg" className="gap-2">
                  Browse Products
                  <ArrowRight className="w-4 h-4" />
                </Button>
              </Link>
              <Link href="/register">
                <Button size="lg" variant="outline">
                  Get Started
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Products Section */}
      <FeaturedProducts />
      
      {/* Features Section */}
      <section className="py-20 bg-surface relative">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <h2 className="text-3xl font-bold text-center mb-12 text-gradient">
            Why Choose Us?
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* Feature 1 */}
            <div className="glass-card p-8 text-center hover:bg-white/[0.05] transition-all duration-300">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mx-auto mb-6">
                <Cpu className="w-6 h-6 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3 text-slate-100">Instant Loading</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Engineered with next-generation edge routing to guarantee sub-millisecond page transitions and zero-delay lookups.
              </p>
            </div>
            
            {/* Feature 2 */}
            <div className="glass-card p-8 text-center hover:bg-white/[0.05] transition-all duration-300">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mx-auto mb-6">
                <ShieldCheck className="w-6 h-6 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3 text-slate-100">Secure Payments</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Your transactions are protected with industry-standard tokenization and processed securely via Midtrans gateway.
              </p>
            </div>
            
            {/* Feature 3 */}
            <div className="glass-card p-8 text-center hover:bg-white/[0.05] transition-all duration-300">
              <div className="w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center mx-auto mb-6">
                <ShoppingBag className="w-6 h-6 text-accent" />
              </div>
              <h3 className="text-xl font-semibold mb-3 text-slate-100">Curated Goods</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Discover a carefully selected collection of premium goods curated by lifestyle experts, updated constantly.
              </p>
            </div>
          </div>
        </div>
      </section>
      
      {/* CTA Section */}
      <section className="py-20 relative overflow-hidden border-t border-white/[0.06]">
        <div className="absolute inset-0 bg-hero-glow opacity-50" />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
          <h2 className="text-3xl font-bold mb-4 text-gradient">Ready to Start Shopping?</h2>
          <p className="text-slate-400 mb-8 max-w-xl mx-auto">
            Create an account today and experience the world's most modern and secure shopping experience.
          </p>
          <Link href="/register">
            <Button size="lg" className="shadow-lg shadow-accent/20">Create Account</Button>
          </Link>
        </div>
      </section>
    </div>
  );
}
