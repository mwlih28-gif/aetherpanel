import React, { useState, useEffect } from 'react'
import { Plus, MapPin, Edit, Trash2, Globe } from 'lucide-react'
import api from '../../lib/api'

interface Location {
  id: string
  short: string
  long: string
  created_at: string
  updated_at: string
}

export default function LocationsPage() {
  const [locations, setLocations] = useState<Location[]>([])
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editingLocation, setEditingLocation] = useState<Location | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchLocations()
  }, [])

  const fetchLocations = async () => {
    try {
      const response = await api.getLocations() as any
      setLocations(response.data || [])
      setLoading(false)
    } catch (error) {
      console.error('Failed to fetch locations:', error)
      setLoading(false)
    }
  }

  const handleCreateLocation = async (data: { short: string; long: string }) => {
    try {
      await api.createLocation(data)
      fetchLocations()
    } catch (error) {
      console.error('Failed to create location:', error)
      alert('Failed to create location: ' + (error as Error).message)
    }
  }

  const handleUpdateLocation = async (id: string, data: { short: string; long: string }) => {
    try {
      await api.updateLocation(id, data)
      fetchLocations()
    } catch (error) {
      console.error('Failed to update location:', error)
      alert('Failed to update location: ' + (error as Error).message)
    }
  }

  const handleDeleteLocation = async (id: string) => {
    if (!confirm('Are you sure you want to delete this location?')) return
    
    try {
      await api.deleteLocation(id)
      fetchLocations()
    } catch (error) {
      console.error('Failed to delete location:', error)
      alert('Failed to delete location: ' + (error as Error).message)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-muted-foreground">Loading locations...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Locations</h1>
          <p className="text-muted-foreground">Manage server locations for node deployment</p>
        </div>
        <button 
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
        >
          <Plus className="h-4 w-4" />
          Create Location
        </button>
      </div>

      {locations.length === 0 ? (
        <div className="text-center py-12">
          <Globe className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">No locations found</h3>
          <p className="text-muted-foreground mb-4">
            Create your first location to organize nodes by geographic region
          </p>
          <button 
            onClick={() => setShowCreateModal(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 mx-auto"
          >
            <Plus className="h-4 w-4" />
            Create Location
          </button>
        </div>
      ) : (
        <div className="bg-card border rounded-lg overflow-hidden">
          <table className="w-full">
            <thead className="bg-muted/50">
              <tr>
                <th className="text-left p-4 font-medium">Short Code</th>
                <th className="text-left p-4 font-medium">Description</th>
                <th className="text-left p-4 font-medium">Nodes</th>
                <th className="text-left p-4 font-medium">Created</th>
                <th className="text-right p-4 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {locations.map((location) => (
                <tr key={location.id} className="border-t hover:bg-muted/25">
                  <td className="p-4">
                    <div className="flex items-center gap-2">
                      <MapPin className="h-4 w-4 text-blue-500" />
                      <span className="font-mono font-medium">{location.short}</span>
                    </div>
                  </td>
                  <td className="p-4">{location.long}</td>
                  <td className="p-4">
                    <span className="text-muted-foreground">0 nodes</span>
                  </td>
                  <td className="p-4 text-muted-foreground">
                    {new Date(location.created_at).toLocaleDateString()}
                  </td>
                  <td className="p-4">
                    <div className="flex items-center gap-2 justify-end">
                      <button
                        onClick={() => setEditingLocation(location)}
                        className="p-2 hover:bg-muted rounded-lg"
                      >
                        <Edit className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => handleDeleteLocation(location.id)}
                        className="p-2 hover:bg-muted rounded-lg text-red-500"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showCreateModal && (
        <LocationModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateLocation}
        />
      )}

      {editingLocation && (
        <LocationModal
          location={editingLocation}
          onClose={() => setEditingLocation(null)}
          onSubmit={(data) => handleUpdateLocation(editingLocation.id, data)}
        />
      )}
    </div>
  )
}

function LocationModal({ 
  location, 
  onClose, 
  onSubmit 
}: { 
  location?: Location
  onClose: () => void
  onSubmit: (data: { short: string; long: string }) => void
}) {
  const [formData, setFormData] = useState({
    short: location?.short || '',
    long: location?.long || ''
  })
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    try {
      await onSubmit(formData)
      onClose()
    } catch (error) {
      console.error('Failed to save location:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-card border rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-semibold mb-4">
          {location ? 'Edit Location' : 'Create Location'}
        </h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Short Code</label>
            <input
              type="text"
              value={formData.short}
              onChange={(e) => setFormData({ ...formData, short: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg font-mono"
              placeholder="us.nyc.1v13"
              pattern="[a-z0-9.-]+"
              title="Only lowercase letters, numbers, dots, and hyphens allowed"
              maxLength={60}
              required
            />
            <p className="text-xs text-muted-foreground mt-1">
              A short identifier used to distinguish this location from others. Must be between 1 and 60 characters, for example, us.nyc.lv13.
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Description</label>
            <textarea
              value={formData.long}
              onChange={(e) => setFormData({ ...formData, long: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg h-20"
              placeholder="New York City, United States"
              maxLength={191}
              required
            />
            <p className="text-xs text-muted-foreground mt-1">
              A longer description of this location. Must be less than 191 characters.
            </p>
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border rounded-lg hover:bg-muted"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg disabled:opacity-50"
            >
              {loading ? 'Saving...' : location ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
