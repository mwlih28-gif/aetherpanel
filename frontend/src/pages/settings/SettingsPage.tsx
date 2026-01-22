import React from 'react'
import { User, Shield, Bell } from 'lucide-react'

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Settings</h1>
        <p className="text-muted-foreground">Manage your account settings</p>
      </div>

      <div className="grid gap-6">
        <div className="bg-card border border-border rounded-xl p-6">
          <div className="flex items-center gap-3 mb-4">
            <User className="w-5 h-5 text-foreground" />
            <h2 className="font-semibold text-foreground">Profile</h2>
          </div>
          <div className="grid gap-4">
            <div>
              <label className="text-sm text-muted-foreground">Username</label>
              <input type="text" defaultValue="admin" className="w-full mt-1 px-4 py-2 bg-input border border-border rounded-lg text-foreground" />
            </div>
            <div>
              <label className="text-sm text-muted-foreground">Email</label>
              <input type="email" defaultValue="admin@example.com" className="w-full mt-1 px-4 py-2 bg-input border border-border rounded-lg text-foreground" />
            </div>
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6">
          <div className="flex items-center gap-3 mb-4">
            <Shield className="w-5 h-5 text-foreground" />
            <h2 className="font-semibold text-foreground">Security</h2>
          </div>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-foreground">Two-Factor Authentication</p>
                <p className="text-sm text-muted-foreground">Add an extra layer of security</p>
              </div>
              <button className="px-4 py-2 bg-primary text-primary-foreground rounded-lg">Enable</button>
            </div>
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6">
          <div className="flex items-center gap-3 mb-4">
            <Bell className="w-5 h-5 text-foreground" />
            <h2 className="font-semibold text-foreground">Notifications</h2>
          </div>
          <div className="space-y-3">
            <label className="flex items-center gap-3 cursor-pointer">
              <input type="checkbox" defaultChecked className="w-4 h-4" />
              <span className="text-foreground">Email notifications</span>
            </label>
            <label className="flex items-center gap-3 cursor-pointer">
              <input type="checkbox" defaultChecked className="w-4 h-4" />
              <span className="text-foreground">Server alerts</span>
            </label>
          </div>
        </div>
      </div>

      <button className="px-6 py-2 bg-primary text-primary-foreground rounded-lg">Save Changes</button>
    </div>
  )
}
