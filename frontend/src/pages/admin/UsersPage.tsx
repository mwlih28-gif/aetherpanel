import React from 'react'
import { Users, Plus } from 'lucide-react'

export default function UsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Users</h1>
          <p className="text-muted-foreground">Manage user accounts</p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg">
          <Plus className="w-4 h-4" /> Add User
        </button>
      </div>
      <div className="bg-card border border-border rounded-xl overflow-hidden">
        <table className="w-full">
          <thead className="bg-accent">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">User</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Role</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {[1, 2, 3].map((i) => (
              <tr key={i} className="hover:bg-accent/50">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
                      <Users className="w-4 h-4 text-primary-foreground" />
                    </div>
                    <span className="text-foreground">user{i}@example.com</span>
                  </div>
                </td>
                <td className="px-6 py-4 text-foreground">User</td>
                <td className="px-6 py-4">
                  <span className="px-2 py-1 bg-green-500/20 text-green-500 rounded text-xs">Active</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
