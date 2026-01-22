const API_BASE_URL = window.location.hostname === 'localhost'
  ? 'http://localhost:8081/api/v1'
  : 'http://178.208.187.30:8081/api/v1'

class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.token = localStorage.getItem('auth_token')
  }

  setToken(token: string) {
    this.token = token
    localStorage.setItem('auth_token', token)
  }

  clearToken() {
    this.token = null
    localStorage.removeItem('auth_token')
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }

    if (options.headers) {
      Object.assign(headers, options.headers)
    }

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Network error' }))
      throw new Error(error.error || `HTTP ${response.status}`)
    }

    return response.json()
  }

  // Auth methods
  async login(credentials: { email: string; password: string }) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    })
  }

  async register(userData: { username: string; email: string; password: string }) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    })
  }

  async logout() {
    return this.request('/auth/logout', { method: 'POST' })
  }

  async getMe() {
    return this.request('/auth/me')
  }

  // Location methods
  async getLocations() {
    return this.request('/locations')
  }

  async createLocation(data: { short: string; long: string }) {
    return this.request('/locations', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateLocation(id: string, data: { short: string; long: string }) {
    return this.request(`/locations/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteLocation(id: string) {
    return this.request(`/locations/${id}`, { method: 'DELETE' })
  }

  // Node methods
  async getNodes() {
    return this.request('/nodes')
  }

  async createNode(data: any) {
    return this.request('/nodes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getNode(id: string) {
    return this.request(`/nodes/${id}`)
  }

  async updateNode(id: string, data: any) {
    return this.request(`/nodes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteNode(id: string) {
    return this.request(`/nodes/${id}`, { method: 'DELETE' })
  }

  async getNodeConfiguration(id: string) {
    return this.request(`/nodes/${id}/configuration`)
  }

  // Server methods
  async getServers() {
    return this.request('/servers')
  }

  async createServer(data: any) {
    return this.request('/servers', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getServer(id: string) {
    return this.request(`/servers/${id}`)
  }

  async updateServer(id: string, data: any) {
    return this.request(`/servers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteServer(id: string) {
    return this.request(`/servers/${id}`, { method: 'DELETE' })
  }

  async startServer(id: string) {
    return this.request(`/servers/${id}/start`, { method: 'POST' })
  }

  async stopServer(id: string) {
    return this.request(`/servers/${id}/stop`, { method: 'POST' })
  }

  async restartServer(id: string) {
    return this.request(`/servers/${id}/restart`, { method: 'POST' })
  }
}

export const api = new ApiClient(API_BASE_URL)
export default api
