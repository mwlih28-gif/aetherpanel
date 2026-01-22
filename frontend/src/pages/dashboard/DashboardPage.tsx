import React from 'react'
import { Server, Users, HardDrive, Activity } from 'lucide-react'

const stats = [
  { name: 'Total Servers', value: '12', icon: Server, color: 'bg-blue-500' },
  { name: 'Active Users', value: '48', icon: Users, color: 'bg-green-500' },
  { name: 'Nodes Online', value: '3/3', icon: HardDrive, color: 'bg-purple-500' },
  { name: 'CPU Usage', value: '42%', icon: Activity, color: 'bg-orange-500' },
]

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Dashboard</h1>
        <p className="text-muted-foreground">Welcome back! Here's an overview of your servers.</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <div key={stat.name} className="bg-card border border-border rounded-xl p-6">
            <div className="flex items-center gap-4">
              <div className={`p-3 rounded-lg ${stat.color}`}>
                <stat.icon className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">{stat.name}</p>
                <p className="text-2xl font-bold text-foreground">{stat.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-card border border-border rounded-xl p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Recent Servers</h2>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-accent/50 rounded-lg">
                <div className="flex items-center gap-3">
                  <div className="w-2 h-2 rounded-full bg-green-500" />
                  <span className="text-foreground">Minecraft Server #{i}</span>
                </div>
                <span className="text-sm text-muted-foreground">Running</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">System Events</h2>
          <div className="space-y-3">
            {['Server started', 'Backup completed', 'User logged in'].map((event, i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-accent/50 rounded-lg">
                <span className="text-foreground">{event}</span>
                <span className="text-sm text-muted-foreground">2m ago</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
