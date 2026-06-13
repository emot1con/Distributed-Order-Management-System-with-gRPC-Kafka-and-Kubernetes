'use client';

import React, { useState, Suspense } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'react-hot-toast';
import { useAuthStore } from '@/store';
import { authService } from '@/services/auth-service';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Store, ArrowLeft } from 'lucide-react';
import { GoogleIcon, FacebookIcon, GithubIcon } from '@/components/ui/oauth-icons';

const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type LoginFormData = z.infer<typeof loginSchema>;

function LoginContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const login = useAuthStore((state) => state.login);
  const [isLoading, setIsLoading] = useState(false);
  
  const returnUrl = searchParams.get('redirect') || '/products';
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  
  React.useEffect(() => {
    if (isAuthenticated) {
      router.replace(returnUrl);
    }
  }, [isAuthenticated, router, returnUrl]);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });
  
  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    try {
      const response = await authService.login(data);
      // Create user object from login data
      const user = {
        id: 0,
        full_name: data.email.split('@')[0], // Temporary fallback
        email: data.email,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      login(response, user);
      toast.success('Successfully logged in!');
      router.push(returnUrl);
    } catch (error: unknown) {
      const errMsg = error instanceof Error ? error.message : 'Login failed. Please try again.';
      toast.error(errMsg);
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleOAuthLogin = (provider: string) => {
    toast.success(`Redirecting to ${provider} login...`);
    let oauthUrl = '';
    if (provider === 'Google') oauthUrl = authService.getGoogleOAuthUrl();
    if (provider === 'Facebook') oauthUrl = authService.getFacebookOAuthUrl();
    if (provider === 'GitHub') oauthUrl = authService.getGithubOAuthUrl();
    
    if (oauthUrl) {
      window.location.href = oauthUrl;
    }
  };
  
  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center bg-surface relative overflow-hidden bg-grid py-12 px-4 sm:px-6 lg:px-8">
      {/* Background glow */}
      <div className="absolute inset-0 bg-hero-glow opacity-60 pointer-events-none" />
      
      <div className="max-w-md w-full relative z-10 animate-slide-up">
        {/* Back Link */}
        <Link href="/" className="inline-flex items-center text-sm text-slate-400 hover:text-slate-200 mb-6 transition-colors">
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Home
        </Link>
        
        {/* Glass Card */}
        <div className="glass-card p-8 shadow-2xl border border-white/[0.06]">
          {/* Logo & Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-accent/15 mb-4">
              <Store className="w-6 h-6 text-accent" />
            </div>
            <h2 className="text-2xl font-bold text-gradient">
              Welcome back
            </h2>
            <p className="text-sm text-slate-400 mt-2">
              Sign in to manage your orders and profile
            </p>
          </div>
          
          {/* Form */}
          <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
            <Input
              label="Email Address"
              type="email"
              placeholder="name@example.com"
              error={errors.email?.message}
              {...register('email')}
            />
            
            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              error={errors.password?.message}
              {...register('password')}
            />
            
            <Button
              type="submit"
              className="w-full mt-2"
              isLoading={isLoading}
            >
              Sign In
            </Button>
          </form>
          
          {/* Divider */}
          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-white/[0.06]"></div>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-surface-card px-2 text-slate-500">Or continue with</span>
            </div>
          </div>
          
          {/* Social Logins */}
          <div className="grid grid-cols-3 gap-3">
            <button
              type="button"
              onClick={() => handleOAuthLogin('Google')}
              className="flex justify-center items-center py-2.5 bg-white/[0.04] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all rounded-lg text-slate-300"
              title="Sign in with Google"
            >
              <GoogleIcon />
            </button>
            <button
              type="button"
              onClick={() => handleOAuthLogin('Facebook')}
              className="flex justify-center items-center py-2.5 bg-white/[0.04] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all rounded-lg text-blue-400 hover:text-blue-300"
              title="Sign in with Facebook"
            >
              <FacebookIcon />
            </button>
            <button
              type="button"
              onClick={() => handleOAuthLogin('GitHub')}
              className="flex justify-center items-center py-2.5 bg-white/[0.04] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all rounded-lg text-slate-200 hover:text-white"
              title="Sign in with GitHub"
            >
              <GithubIcon />
            </button>
          </div>
          
          {/* Footer Link */}
          <p className="text-center text-sm text-slate-400 mt-8">
            Don't have an account?{' '}
            <Link href="/register" className="text-accent hover:text-accent-hover font-semibold transition-colors">
              Create an account
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-accent" />
      </div>
    }>
      <LoginContent />
    </Suspense>
  );
}
