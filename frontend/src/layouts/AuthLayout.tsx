import { Outlet } from 'react-router-dom'

export default function AuthLayout() {
  return (
    <div className="min-h-screen flex dark bg-background">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-red-600 to-orange-500 p-12 flex-col justify-between">
        <div>
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-white/20 backdrop-blur rounded-xl flex items-center justify-center">
              <span className="text-white font-bold text-2xl">A</span>
            </div>
            <span className="text-3xl font-bold text-white">Aether Panel</span>
          </div>
        </div>
        
        <div className="space-y-6">
          <h1 className="text-4xl font-bold text-white leading-tight">
            Next-Generation<br />Game Server Management
          </h1>
          <p className="text-white/80 text-lg max-w-md">
            Deploy, manage, and scale your game servers with ease. 
            Built for performance, designed for simplicity.
          </p>
          <div className="flex gap-8 text-white/90">
            <div>
              <div className="text-3xl font-bold">50+</div>
              <div className="text-sm text-white/70">Supported Games</div>
            </div>
            <div>
              <div className="text-3xl font-bold">99.9%</div>
              <div className="text-sm text-white/70">Uptime SLA</div>
            </div>
            <div>
              <div className="text-3xl font-bold">24/7</div>
              <div className="text-sm text-white/70">Support</div>
            </div>
          </div>
        </div>

        <div className="text-white/60 text-sm">
          © 2024 Aether Panel. All rights reserved.
        </div>
      </div>

      {/* Right side - Auth form */}
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="w-full max-w-md">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
